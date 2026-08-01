package danfse

import (
	"bytes"
	_ "embed"
	"fmt"
	"image/png"
	"strconv"
	"strings"

	"github.com/signintech/gopdf"
	"github.com/skip2/go-qrcode"
)

//go:embed fonts/LiberationSans-Regular.ttf
var embeddedFontReg []byte

//go:embed fonts/LiberationSans-Bold.ttf
var embeddedFontBold []byte

//go:embed assets/logo-nfse.png
var embeddedLogo []byte

type DanfseRenderer interface {
	Render(xmlData []byte) ([]byte, error)
}

type renderer struct {
	fontRegData  []byte
	fontBoldData []byte
	logoData     []byte
}

func NewDanfseRenderer() DanfseRenderer {
	return &renderer{
		fontRegData:  embeddedFontReg,
		fontBoldData: embeddedFontBold,
		logoData:     embeddedLogo,
	}
}

func cmToPt(cm float64) float64 {
	return cm * 28.34645669
}

func formatDoc(doc string) string {
	if len(doc) == 14 {
		return fmt.Sprintf("%s.%s.%s/%s-%s", doc[0:2], doc[2:5], doc[5:8], doc[8:12], doc[12:14])
	} else if len(doc) == 11 {
		return fmt.Sprintf("%s.%s.%s-%s", doc[0:3], doc[3:6], doc[6:9], doc[9:11])
	}
	return doc
}

func formatCEP(cep string) string {
	if len(cep) == 8 {
		return fmt.Sprintf("%s.%s-%s", cep[0:2], cep[2:5], cep[5:8])
	}
	return cep
}

func formatTribISSQN(code string) string {
	switch code {
	case "1":
		return "Operação Tributável"
	case "2":
		return "Exportação de Serviço"
	case "3":
		return "Não Incidência"
	case "4":
		return "Imunidade"
	case "":
		return "-"
	default:
		return code
	}
}

func formatRegEspTrib(code string) string {
	switch code {
	case "0":
		return "Nenhum"
	case "1":
		return "Microempresa Municipal"
	case "2":
		return "Estimativa"
	case "3":
		return "Sociedade de Profissionais"
	case "4":
		return "Cooperativa"
	case "5":
		return "MEI"
	case "6":
		return "ME/EPP"
	case "9":
		return "Outros"
	case "":
		return "-"
	default:
		return code
	}
}

func formatTpRetISSQN(code string) string {
	switch code {
	case "1":
		return "Não Retido"
	case "2":
		return "Retido pelo Tomador"
	case "3":
		return "Retido pelo Intermediário"
	case "":
		return "-"
	default:
		return code
	}
}

func formatTpSusp(code string) string {
	switch code {
	case "0":
		return "Não"
	case "1":
		return "Exigibilidade Suspensa por Decisão Judicial"
	case "2":
		return "Exigibilidade Suspensa por Processo Administrativo"
	case "":
		return "-"
	default:
		return code
	}
}

func formatCurrency(val string) string {
	if val == "" || val == "-" {
		return val
	}
	f, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return val
	}
	parts := strings.Split(fmt.Sprintf("%.2f", f), ".")
	intPart := parts[0]
	decPart := parts[1]
	var withSeps string
	for i, c := range intPart {
		if i > 0 && (len(intPart)-i)%3 == 0 {
			withSeps += "."
		}
		withSeps += string(c)
	}
	return "R$ " + withSeps + "," + decPart
}

// ufByIBGEPrefix mapeia os dois primeiros dígitos do código IBGE do município
// para a sigla da UF. Usado para derivar a UF quando o XML não a informa
// (ex.: <endNac> do tomador só traz cMun e CEP, sem UF).
var ufByIBGEPrefix = map[string]string{
	"11": "RO", "12": "AC", "13": "AM", "14": "RR", "15": "PA", "16": "AP", "17": "TO",
	"21": "MA", "22": "PI", "23": "CE", "24": "RN", "25": "PB", "26": "PE", "27": "AL", "28": "SE", "29": "BA",
	"31": "MG", "32": "ES", "33": "RJ", "35": "SP",
	"41": "PR", "42": "SC", "43": "RS",
	"50": "MS", "51": "MT", "52": "GO", "53": "DF",
}

// ufFromIBGE deriva a sigla da UF a partir do código IBGE do município (7 dígitos);
// retorna "" quando o código é inválido ou desconhecido.
func ufFromIBGE(cMun string) string {
	if len(cMun) < 2 {
		return ""
	}
	return ufByIBGEPrefix[cMun[:2]]
}

func formatMunUF(xMun, cMun, uf string) string {
	mun := xMun
	if mun == "" {
		mun = ibgeCities[cMun]
	}
	if mun == "" {
		mun = cMun
	}
	if uf == "" {
		// O <endNac> do tomador não traz UF; deriva-se do código IBGE.
		uf = ufFromIBGE(cMun)
	}
	if mun == "" {
		if uf != "" {
			return uf
		}
		return "-"
	}
	if uf != "" {
		return mun + " / " + uf
	}
	return mun
}

func formatOpSimpNac(code string) string {
	switch code {
	case "1":
		return "1 - Não Optante"
	case "2":
		return "2 - Optante - MEI"
	case "3":
		return "3 - Optante - ME/EPP"
	case "":
		return "-"
	default:
		return code
	}
}

func formatRegApTribSN(code string) string {
	switch code {
	case "1":
		return "1 - ME/EPP"
	case "2":
		return "2 - MEI"
	case "":
		return "-"
	default:
		return code
	}
}

func formatTpRetPisCofins(code string) string {
	switch code {
	case "1":
		return "PIS/COFINS/CSLL Não Retido"
	case "2":
		return "Retido pelo Tomador"
	case "":
		return "-"
	default:
		return code
	}
}

func formatTribNac(code string) string {
	if len(code) == 6 {
		return code[0:2] + "." + code[2:4] + "." + code[4:6]
	}
	return code
}

func formatPhone(p string) string {
	if len(p) == 10 {
		return fmt.Sprintf("(%s) %s-%s", p[0:2], p[2:6], p[6:10])
	} else if len(p) == 11 {
		return fmt.Sprintf("(%s) %s-%s", p[0:2], p[2:7], p[7:11])
	}
	return p
}

func (r *renderer) Render(xmlData []byte) ([]byte, error) {
	nfse, err := ParseXML(xmlData)
	if err != nil {
		return nil, fmt.Errorf("pdf: falha ao ler xml: %w", err)
	}

	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	pdf.AddPage()

	err = pdf.AddTTFFontByReader("LiberationSans", bytes.NewReader(r.fontRegData))
	if err != nil {
		return nil, fmt.Errorf("pdf: erro fonte regular: %w", err)
	}
	err = pdf.AddTTFFontByReader("LiberationSans-Bold", bytes.NewReader(r.fontBoldData))
	if err != nil {
		return nil, fmt.Errorf("pdf: erro fonte bold: %w", err)
	}

	// Borda da Página (1 pt)
	pdf.SetLineWidth(1.0)
	pdf.RectFromUpperLeftWithStyle(cmToPt(0.15), cmToPt(0.15), cmToPt(20.70), cmToPt(29.40), "D")

	// Linhas internas (0.5 pt)
	pdf.SetLineWidth(0.5)

	// HEADER
	// Header Cinza (5%)
	pdf.SetFillColor(242, 242, 242)
	pdf.RectFromUpperLeftWithStyle(cmToPt(0.30), cmToPt(0.30), cmToPt(20.40), cmToPt(1.17), "F")
	pdf.SetFillColor(0, 0, 0)

	// Logo (aspect ratio fix)
	var logoW, logoH float64 = cmToPt(4.00), cmToPt(1.16)
	if img, decErr := png.DecodeConfig(bytes.NewReader(r.logoData)); decErr == nil {
		ratio := float64(img.Height) / float64(img.Width)
		logoH = logoW * ratio
	}
	logoY := cmToPt(0.30)
	if logoH < cmToPt(1.16) {
		logoY += (cmToPt(1.16) - logoH) / 2
	}
	logoHolder, err := gopdf.ImageHolderByBytes(r.logoData)
	if err != nil {
		return nil, err
	}
	err = pdf.ImageByHolder(logoHolder, cmToPt(0.49), logoY, &gopdf.Rect{W: logoW, H: logoH})
	if err != nil {
		return nil, err
	}

	// Textos Centrais
	pdf.SetFont("LiberationSans-Bold", "", 9)
	pdf.SetX(cmToPt(5.41))
	pdf.SetY(cmToPt(0.35))
	pdf.CellWithOption(&gopdf.Rect{W: cmToPt(10.19), H: cmToPt(0.35)}, "DANFSe v2.0", gopdf.CellOption{Align: gopdf.Center | gopdf.Top})

	pdf.SetY(cmToPt(0.65))
	pdf.SetX(cmToPt(5.41))
	pdf.CellWithOption(&gopdf.Rect{W: cmToPt(10.19), H: cmToPt(0.35)}, "Documento Auxiliar da NFS-e", gopdf.CellOption{Align: gopdf.Center | gopdf.Top})

	// Homologação: o aviso "SEM VALIDADE JURÍDICA" depende do TIPO DE AMBIENTE
	// (tpAmb: 1=produção, 2=homologação), NÃO do AmbGer (1=Prefeitura, 2=SEFIN
	// Nacional, que é 2 também em produção).
	if nfse.InfNFSe.DPS.InfDPS.TpAmb == "2" {
		pdf.SetTextColor(255, 0, 0)
		pdf.SetY(cmToPt(0.95))
		pdf.SetX(cmToPt(5.41))
		pdf.CellWithOption(&gopdf.Rect{W: cmToPt(10.19), H: cmToPt(0.35)}, "NFS-e SEM VALIDADE JURÍDICA", gopdf.CellOption{Align: gopdf.Center | gopdf.Top})
		pdf.SetTextColor(0, 0, 0)
	}

	// Textos Direita
	pdf.SetFont("LiberationSans", "", 8)
	pdf.SetY(cmToPt(0.45))
	pdf.SetX(cmToPt(15.62))
	uf := nfse.InfNFSe.Emit.EnderNac.UF
	if uf == "" {
		uf = ufFromIBGE(nfse.InfNFSe.Emit.EnderNac.CMun)
	}
	if uf == "" {
		return nil, fmt.Errorf("pdf: UF do emitente ausente e não derivável do cMun %q", nfse.InfNFSe.Emit.EnderNac.CMun)
	}
	// Separador " / " conforme NT 008/2026 v1.02, tabela 2.4.5 ("Município: CCCC / CC").
	pdf.Cell(nil, fmt.Sprintf("Município: %s / %s", nfse.InfNFSe.XLocEmi, uf))

	pdf.SetFont("LiberationSans", "", 6)
	pdf.SetY(cmToPt(0.85))
	pdf.SetX(cmToPt(15.62))
	pdf.Cell(nil, "Ambiente Gerador: "+nfse.InfNFSe.AmbGer)

	pdf.SetY(cmToPt(1.15))
	pdf.SetX(cmToPt(15.62))
	pdf.Cell(nil, "Tipo de Ambiente: "+nfse.InfNFSe.DPS.InfDPS.TpAmb)

	// LINHA SEPARADORA DO HEADER
	pdf.Line(cmToPt(0.30), cmToPt(1.47), cmToPt(20.70), cmToPt(1.47))

	// CHAVE DE ACESSO (Sem fundo cinza)
	pdf.SetFont("LiberationSans-Bold", "", 7)
	pdf.SetY(cmToPt(1.55))
	pdf.SetX(cmToPt(0.40))
	pdf.Cell(nil, "CHAVE DE ACESSO DA NFS-E")

	pdf.SetFont("LiberationSans", "", 7)
	pdf.SetY(cmToPt(1.90))
	pdf.SetX(cmToPt(0.40))
	chaveAcessoTop := strings.TrimPrefix(nfse.InfNFSe.ID, "NFS")
	pdf.Cell(nil, chaveAcessoTop)

	// QR Code (Sem borda)
	qrUrl := "https://www.nfse.gov.br/ConsultaPublica/?tpc=1&chave=" + chaveAcessoTop
	// Sem o QR não há como verificar a autenticidade da NFS-e: qualquer falha da
	// cadeia aborta a geração em vez de emitir um DANFSe incompleto.
	q, err := qrcode.New(qrUrl, qrcode.Medium)
	if err != nil {
		return nil, fmt.Errorf("pdf: qrcode: %w", err)
	}
	q.DisableBorder = true
	qrBytes, err := q.PNG(256)
	if err != nil {
		return nil, fmt.Errorf("pdf: qrcode: %w", err)
	}
	imgH, err := gopdf.ImageHolderByBytes(qrBytes)
	if err != nil {
		return nil, fmt.Errorf("pdf: qrcode: %w", err)
	}
	// Y ajustado para 1.50 (encosta levemente mais no topo para dar respiro embaixo)
	err = pdf.ImageByHolder(imgH, cmToPt(17.48), cmToPt(1.60), &gopdf.Rect{W: cmToPt(1.90), H: cmToPt(1.90)})
	if err != nil {
		return nil, fmt.Errorf("pdf: qrcode: %w", err)
	}

	// A autenticidade...
	pdf.SetFont("LiberationSans", "", 6)
	// Textos descem mais 0.10cm para afastar do QR Code (ficam em 3.65, 3.85 e 4.05)
	pdf.SetY(cmToPt(3.65))
	pdf.SetX(cmToPt(15.62))
	pdf.Cell(nil, "A autenticidade desta NFS-e pode ser verificada")

	pdf.SetY(cmToPt(3.85))
	pdf.SetX(cmToPt(15.62))
	pdf.Cell(nil, "pela leitura deste código QR ou pela consulta da")

	pdf.SetY(cmToPt(4.05)) // Base fica em ~4.26, mantendo margem da linha 4.34
	pdf.SetX(cmToPt(15.62))
	pdf.Cell(nil, "chave de acesso no portal nacional da NFS-e")

	// DADOS DA NOTA (NÚMERO, COMPETÊNCIA, EMISSÃO)
	pdf.SetFont("LiberationSans-Bold", "", 7)
	pdf.SetY(cmToPt(2.35))
	pdf.SetX(cmToPt(0.40))
	pdf.Cell(nil, "NÚMERO DA NFS-E")

	pdf.SetX(cmToPt(5.41))
	pdf.Cell(nil, "COMPETÊNCIA DA NFS-E")

	pdf.SetX(cmToPt(10.51))
	pdf.Cell(nil, "DATA E HORA DA EMISSÃO DA NFS-E")

	pdf.SetFont("LiberationSans", "", 7)
	pdf.SetY(cmToPt(2.65))
	pdf.SetX(cmToPt(0.40))
	pdf.Cell(nil, nfse.InfNFSe.NNFSe)

	pdf.SetX(cmToPt(5.41))
	pdf.Cell(nil, nfse.InfNFSe.DPS.InfDPS.DCompet)

	pdf.SetX(cmToPt(10.51))
	pdf.Cell(nil, nfse.InfNFSe.DhProc)

	// DADOS DA DPS
	pdf.SetFont("LiberationSans-Bold", "", 7)
	pdf.SetY(cmToPt(3.05))
	pdf.SetX(cmToPt(0.40))
	pdf.Cell(nil, "NÚMERO DA DPS")

	pdf.SetX(cmToPt(5.41))
	pdf.Cell(nil, "SÉRIE DA DPS")

	pdf.SetX(cmToPt(10.51))
	pdf.Cell(nil, "DATA E HORA DA EMISSÃO DA DPS")

	pdf.SetFont("LiberationSans", "", 7)
	pdf.SetY(cmToPt(3.35))
	pdf.SetX(cmToPt(0.40))
	pdf.Cell(nil, nfse.InfNFSe.DPS.InfDPS.NDPS)

	pdf.SetX(cmToPt(5.41))
	pdf.Cell(nil, nfse.InfNFSe.DPS.InfDPS.Serie)

	pdf.SetX(cmToPt(10.51))
	pdf.Cell(nil, nfse.InfNFSe.DPS.InfDPS.DhEmi)

	// EMITENTE DA NFS-e / SITUAÇÃO DA NFS-e / FINALIDADE
	// Fundo cinza APENAS para a célula "EMITENTE DA NFS-e"
	pdf.SetFillColor(242, 242, 242)
	pdf.RectFromUpperLeftWithStyle(cmToPt(0.30), cmToPt(3.67), cmToPt(4.90), cmToPt(0.67), "F")
	pdf.SetFillColor(0, 0, 0)

	pdf.SetFont("LiberationSans-Bold", "", 7)
	pdf.SetY(cmToPt(3.72))
	pdf.SetX(cmToPt(0.40))
	pdf.Cell(nil, "EMITENTE DA NFS-e")

	pdf.SetX(cmToPt(5.41))
	pdf.Cell(nil, "SITUAÇÃO DA NFS-e")

	pdf.SetX(cmToPt(10.51))
	pdf.Cell(nil, "FINALIDADE")

	pdf.SetFont("LiberationSans", "", 7)
	pdf.SetY(cmToPt(4.02))
	pdf.SetX(cmToPt(5.41))
	pdf.Cell(nil, "1 - Normal") // Fixo no exemplo

	pdf.SetX(cmToPt(10.51))
	pdf.Cell(nil, "1 - Normal") // Fixo no exemplo

	// PRESTADOR / FORNECEDOR
	// Fundo cinza para a célula de título (Y=4.34, alt=0.63)
	pdf.SetFillColor(242, 242, 242)
	pdf.RectFromUpperLeftWithStyle(cmToPt(0.30), cmToPt(4.34), cmToPt(4.90), cmToPt(0.64), "F")
	pdf.SetFillColor(0, 0, 0)

	// LINHA DIVISÓRIA ABAIXO DO EMITENTE (Cruza a página toda, desenhada depois do fundo para sobrepor)
	pdf.Line(cmToPt(0.30), cmToPt(4.34), cmToPt(20.70), cmToPt(4.34))

	pdf.SetFont("LiberationSans-Bold", "", 7)
	pdf.SetY(cmToPt(4.41))
	pdf.SetX(cmToPt(0.40))
	pdf.Cell(nil, "PRESTADOR / FORNECEDOR")

	pdf.SetFont("LiberationSans-Bold", "", 6)
	pdf.SetX(cmToPt(5.41))
	pdf.Cell(nil, "CNPJ / CPF / NIF")

	pdf.SetX(cmToPt(10.51))
	pdf.Cell(nil, "Indicador Municipal (Inscrição)")

	pdf.SetX(cmToPt(15.62))
	pdf.Cell(nil, "Telefone")

	pdf.SetFont("LiberationSans", "", 7)
	pdf.SetY(cmToPt(4.71))
	pdf.SetX(cmToPt(5.41))
	pdf.Cell(nil, formatDoc(nfse.InfNFSe.Emit.CNPJ))

	pdf.SetX(cmToPt(10.51))
	pdf.Cell(nil, nfse.InfNFSe.DPS.InfDPS.Prest.IM)

	fone := nfse.InfNFSe.Emit.Fone
	if fone == "" {
		fone = nfse.InfNFSe.DPS.InfDPS.Prest.Fone
	}
	pdf.SetX(cmToPt(15.62))
	pdf.Cell(nil, formatPhone(fone))

	// Linha 2 do Prestador
	pdf.SetFont("LiberationSans-Bold", "", 6)
	pdf.SetY(cmToPt(5.05))
	pdf.SetX(cmToPt(0.40))
	pdf.Cell(nil, "Nome / Nome Empresarial")

	pdf.SetX(cmToPt(10.51))
	pdf.Cell(nil, "Município / Sigla UF")

	pdf.SetX(cmToPt(15.62))
	pdf.Cell(nil, "Código IBGE / CEP")

	pdf.SetFont("LiberationSans", "", 7)
	pdf.SetY(cmToPt(5.35))
	pdf.SetX(cmToPt(0.40))
	pdf.Cell(nil, nfse.InfNFSe.Emit.XNome)

	munUF := formatMunUF(nfse.InfNFSe.Emit.EnderNac.XMun, nfse.InfNFSe.Emit.EnderNac.CMun, nfse.InfNFSe.Emit.EnderNac.UF)
	pdf.SetX(cmToPt(10.51))
	pdf.Cell(nil, munUF)

	ibgeCep := "-"
	if nfse.InfNFSe.Emit.EnderNac.CEP != "" {
		ibgeCep = nfse.InfNFSe.Emit.EnderNac.CMun + " / " + formatCEP(nfse.InfNFSe.Emit.EnderNac.CEP)
	}
	pdf.SetX(cmToPt(15.62))
	pdf.Cell(nil, ibgeCep)

	// Linha 3 do Prestador
	pdf.SetFont("LiberationSans-Bold", "", 6)
	pdf.SetY(cmToPt(5.69))
	pdf.SetX(cmToPt(0.40))
	pdf.Cell(nil, "Endereço")

	pdf.SetX(cmToPt(10.51))
	pdf.Cell(nil, "E-mail")

	pdf.SetFont("LiberationSans", "", 7)
	pdf.SetY(cmToPt(5.99))
	pdf.SetX(cmToPt(0.40))
	end := nfse.InfNFSe.Emit.EnderNac.XLgr
	if nfse.InfNFSe.Emit.EnderNac.Nro != "" {
		end += ", " + nfse.InfNFSe.Emit.EnderNac.Nro
	}
	if nfse.InfNFSe.Emit.EnderNac.XCpl != "" {
		end += ", " + nfse.InfNFSe.Emit.EnderNac.XCpl
	}
	if nfse.InfNFSe.Emit.EnderNac.XBairro != "" {
		end += ", " + nfse.InfNFSe.Emit.EnderNac.XBairro
	}
	pdf.Cell(nil, end)

	email := nfse.InfNFSe.Emit.Email
	if email == "" {
		email = nfse.InfNFSe.DPS.InfDPS.Prest.Email
	}
	pdf.SetX(cmToPt(10.51))
	pdf.Cell(nil, email)

	// Linha 4 do Prestador
	pdf.SetFont("LiberationSans-Bold", "", 6)
	pdf.SetY(cmToPt(6.35))
	pdf.SetX(cmToPt(0.40))
	pdf.Cell(nil, "Simples Nacional na Data de Competência")

	pdf.SetX(cmToPt(10.51))
	pdf.Cell(nil, "Regime de Apuração Tributária pelo SN")

	pdf.SetFont("LiberationSans", "", 7)
	pdf.SetY(cmToPt(6.65))
	pdf.SetX(cmToPt(0.40))
	opSimp := formatOpSimpNac(nfse.InfNFSe.DPS.InfDPS.Prest.RegTrib.OpSimpNac)
	pdf.Cell(nil, opSimp)

	pdf.SetX(cmToPt(10.51))
	regAp := formatRegApTribSN(nfse.InfNFSe.DPS.InfDPS.Prest.RegTrib.RegApTribSN)
	pdf.Cell(nil, regAp)

	// TOMADOR / ADQUIRENTE
	pdf.SetFillColor(242, 242, 242)
	pdf.RectFromUpperLeftWithStyle(cmToPt(0.30), cmToPt(6.92), cmToPt(4.90), cmToPt(0.64), "F")
	pdf.SetFillColor(0, 0, 0)

	// LINHA DIVISÓRIA ABAIXO DO PRESTADOR (desenhada depois do fundo para sobrepor)
	pdf.Line(cmToPt(0.30), cmToPt(6.92), cmToPt(20.70), cmToPt(6.92))

	pdf.SetFont("LiberationSans-Bold", "", 7)
	pdf.SetY(cmToPt(6.99))
	pdf.SetX(cmToPt(0.40))
	pdf.Cell(nil, "TOMADOR / ADQUIRENTE")

	toma := nfse.InfNFSe.DPS.InfDPS.Toma
	if toma.XNome == "" && toma.CNPJ == "" {
		// Tomador não identificado
		pdf.SetFont("LiberationSans", "", 7)
		pdf.SetY(cmToPt(7.29))
		pdf.SetX(cmToPt(5.41))
		pdf.Cell(nil, "TOMADOR/ADQUIRENTE DA OPERAÇÃO NÃO IDENTIFICADO NA NFS-e")
	} else {
		pdf.SetFont("LiberationSans-Bold", "", 6)
		pdf.SetX(cmToPt(5.41))
		pdf.Cell(nil, "CNPJ / CPF / NIF")

		pdf.SetX(cmToPt(10.51))
		pdf.Cell(nil, "Indicador Municipal (Inscrição)")

		pdf.SetX(cmToPt(15.62))
		pdf.Cell(nil, "Telefone")

		pdf.SetFont("LiberationSans", "", 7)
		pdf.SetY(cmToPt(7.29))
		pdf.SetX(cmToPt(5.41))
		pdf.Cell(nil, formatDoc(toma.CNPJ))

		pdf.SetX(cmToPt(10.51))
		pdf.Cell(nil, toma.IM)

		pdf.SetX(cmToPt(15.62))
		pdf.Cell(nil, formatPhone(toma.Fone))

		// Linha 2 do Tomador
		pdf.SetFont("LiberationSans-Bold", "", 6)
		pdf.SetY(cmToPt(7.63))
		pdf.SetX(cmToPt(0.40))
		pdf.Cell(nil, "Nome / Nome Empresarial")

		pdf.SetX(cmToPt(10.51))
		pdf.Cell(nil, "Município / Sigla UF")

		pdf.SetX(cmToPt(15.62))
		pdf.Cell(nil, "Código IBGE / CEP")

		pdf.SetFont("LiberationSans", "", 7)
		pdf.SetY(cmToPt(7.93))
		pdf.SetX(cmToPt(0.40))
		pdf.Cell(nil, toma.XNome)

		munUFToma := formatMunUF(toma.End.EndNac.XMun, toma.End.EndNac.CMun, toma.End.EndNac.UF)
		pdf.SetX(cmToPt(10.51))
		pdf.Cell(nil, munUFToma)

		ibgeCepToma := "-"
		if toma.End.EndNac.CEP != "" {
			ibgeCepToma = toma.End.EndNac.CMun + " / " + formatCEP(toma.End.EndNac.CEP)
		}
		pdf.SetX(cmToPt(15.62))
		pdf.Cell(nil, ibgeCepToma)

		// Linha 3 do Tomador
		pdf.SetFont("LiberationSans-Bold", "", 6)
		pdf.SetY(cmToPt(8.29))
		pdf.SetX(cmToPt(0.40))
		pdf.Cell(nil, "Endereço")

		pdf.SetX(cmToPt(10.51))
		pdf.Cell(nil, "E-mail")

		pdf.SetFont("LiberationSans", "", 7)
		pdf.SetY(cmToPt(8.59))
		pdf.SetX(cmToPt(0.40))
		endToma := toma.End.XLgr
		if toma.End.Nro != "" {
			endToma += ", " + toma.End.Nro
		}
		if toma.End.XCpl != "" {
			endToma += ", " + toma.End.XCpl
		}
		if toma.End.XBairro != "" {
			endToma += ", " + toma.End.XBairro
		}
		pdf.Cell(nil, endToma)

		pdf.SetX(cmToPt(10.51))
		pdf.Cell(nil, toma.Email)
	}

	// DESTINATÁRIO DA OPERAÇÃO
	dest := nfse.InfNFSe.DPS.InfDPS.IBSCBS.Dest

	pdf.SetFillColor(242, 242, 242)
	pdf.RectFromUpperLeftWithStyle(cmToPt(0.30), cmToPt(8.86), cmToPt(4.90), cmToPt(0.64), "F")
	pdf.SetFillColor(0, 0, 0)

	// LINHA DIVISÓRIA ABAIXO DO TOMADOR
	pdf.Line(cmToPt(0.30), cmToPt(8.86), cmToPt(20.70), cmToPt(8.86))

	pdf.SetFont("LiberationSans-Bold", "", 7)
	pdf.SetY(cmToPt(8.93))
	pdf.SetX(cmToPt(0.40))
	pdf.Cell(nil, "DESTINATÁRIO DA OPERAÇÃO")

	if dest.XNome == "" && dest.CNPJ == "" {
		pdf.SetFont("LiberationSans", "", 7)
		pdf.SetY(cmToPt(9.23))
		pdf.SetX(cmToPt(5.41))
		pdf.Cell(nil, "DESTINATÁRIO DA OPERAÇÃO NÃO IDENTIFICADO NA NFS-e")
	} else {
		pdf.SetFont("LiberationSans-Bold", "", 6)
		pdf.SetX(cmToPt(5.41))
		pdf.Cell(nil, "CNPJ / CPF / NIF")
		pdf.SetX(cmToPt(15.62))
		pdf.Cell(nil, "Telefone")

		pdf.SetFont("LiberationSans", "", 7)
		pdf.SetY(cmToPt(9.23))
		pdf.SetX(cmToPt(5.41))
		pdf.Cell(nil, formatDoc(dest.CNPJ))

		pdf.SetX(cmToPt(15.62))
		pdf.Cell(nil, formatPhone(dest.Fone))

		// Linha 2
		pdf.SetFont("LiberationSans-Bold", "", 6)
		pdf.SetY(cmToPt(9.57))
		pdf.SetX(cmToPt(0.40))
		pdf.Cell(nil, "Nome / Nome Empresarial")

		pdf.SetX(cmToPt(10.51))
		pdf.Cell(nil, "Município / Sigla UF")

		pdf.SetX(cmToPt(15.62))
		pdf.Cell(nil, "Código IBGE / CEP")

		pdf.SetFont("LiberationSans", "", 7)
		pdf.SetY(cmToPt(9.87))
		pdf.SetX(cmToPt(0.40))
		pdf.Cell(nil, dest.XNome)

		munUFDest := formatMunUF(dest.XMun, dest.CMun, dest.UF)
		pdf.SetX(cmToPt(10.51))
		pdf.Cell(nil, munUFDest)

		ibgeCepDest := "-"
		if dest.CEP != "" {
			ibgeCepDest = dest.CMun + " / " + formatCEP(dest.CEP)
		}
		pdf.SetX(cmToPt(15.62))
		pdf.Cell(nil, ibgeCepDest)

		// Linha 3
		pdf.SetFont("LiberationSans-Bold", "", 6)
		pdf.SetY(cmToPt(10.23))
		pdf.SetX(cmToPt(0.40))
		pdf.Cell(nil, "Endereço")

		pdf.SetX(cmToPt(10.51))
		pdf.Cell(nil, "E-mail")

		pdf.SetFont("LiberationSans", "", 7)
		pdf.SetY(cmToPt(10.53))
		pdf.SetX(cmToPt(0.40))
		endD := dest.End.XLgr
		if dest.End.Nro != "" {
			endD += ", " + dest.End.Nro
		}
		if dest.End.XCpl != "" {
			endD += ", " + dest.End.XCpl
		}
		if dest.End.XBairro != "" {
			endD += ", " + dest.End.XBairro
		}
		pdf.Cell(nil, endD)

		pdf.SetX(cmToPt(10.51))
		pdf.Cell(nil, dest.Email)
	}

	// INTERMEDIÁRIO DA OPERAÇÃO
	interm := nfse.InfNFSe.DPS.InfDPS.Interm

	pdf.SetFillColor(242, 242, 242)
	pdf.RectFromUpperLeftWithStyle(cmToPt(0.30), cmToPt(10.80), cmToPt(4.90), cmToPt(0.64), "F")
	pdf.SetFillColor(0, 0, 0)

	// LINHA DIVISÓRIA ABAIXO DO DESTINATÁRIO
	pdf.Line(cmToPt(0.30), cmToPt(10.80), cmToPt(20.70), cmToPt(10.80))

	pdf.SetFont("LiberationSans-Bold", "", 7)
	pdf.SetY(cmToPt(10.87))
	pdf.SetX(cmToPt(0.40))
	pdf.Cell(nil, "INTERMEDIÁRIO DA OPERAÇÃO")

	if interm.XNome == "" && interm.CNPJ == "" {
		pdf.SetFont("LiberationSans", "", 7)
		pdf.SetY(cmToPt(11.17))
		pdf.SetX(cmToPt(5.41))
		pdf.Cell(nil, "INTERMEDIÁRIO DA OPERAÇÃO NÃO IDENTIFICADO NA NFS-e")
	} else {
		pdf.SetFont("LiberationSans-Bold", "", 6)
		pdf.SetX(cmToPt(5.41))
		pdf.Cell(nil, "CNPJ / CPF / NIF")
		pdf.SetX(cmToPt(10.51))
		pdf.Cell(nil, "Indicador Municipal (Inscrição)")
		pdf.SetX(cmToPt(15.62))
		pdf.Cell(nil, "Telefone")

		pdf.SetFont("LiberationSans", "", 7)
		pdf.SetY(cmToPt(11.17))
		pdf.SetX(cmToPt(5.41))
		pdf.Cell(nil, formatDoc(interm.CNPJ))

		pdf.SetX(cmToPt(10.51))
		pdf.Cell(nil, interm.IM)

		pdf.SetX(cmToPt(15.62))
		pdf.Cell(nil, formatPhone(interm.Fone))

		// Linha 2
		pdf.SetFont("LiberationSans-Bold", "", 6)
		pdf.SetY(cmToPt(11.51))
		pdf.SetX(cmToPt(0.40))
		pdf.Cell(nil, "Nome / Nome Empresarial")

		pdf.SetX(cmToPt(10.51))
		pdf.Cell(nil, "Município / Sigla UF")

		pdf.SetX(cmToPt(15.62))
		pdf.Cell(nil, "Código IBGE / CEP")

		pdf.SetFont("LiberationSans", "", 7)
		pdf.SetY(cmToPt(11.81))
		pdf.SetX(cmToPt(0.40))
		pdf.Cell(nil, interm.XNome)

		munUFInt := formatMunUF(interm.XMun, interm.CMun, interm.UF)
		pdf.SetX(cmToPt(10.51))
		pdf.Cell(nil, munUFInt)

		ibgeCepInt := "-"
		if interm.CEP != "" {
			ibgeCepInt = interm.CMun + " / " + formatCEP(interm.CEP)
		}
		pdf.SetX(cmToPt(15.62))
		pdf.Cell(nil, ibgeCepInt)

		// Linha 3
		pdf.SetFont("LiberationSans-Bold", "", 6)
		pdf.SetY(cmToPt(12.16))
		pdf.SetX(cmToPt(0.40))
		pdf.Cell(nil, "Endereço")

		pdf.SetX(cmToPt(10.51))
		pdf.Cell(nil, "E-mail")

		pdf.SetFont("LiberationSans", "", 7)
		pdf.SetY(cmToPt(12.46))
		pdf.SetX(cmToPt(0.40))
		endI := interm.End.XLgr
		if interm.End.Nro != "" {
			endI += ", " + interm.End.Nro
		}
		if interm.End.XCpl != "" {
			endI += ", " + interm.End.XCpl
		}
		if interm.End.XBairro != "" {
			endI += ", " + interm.End.XBairro
		}
		pdf.Cell(nil, endI)

		pdf.SetX(cmToPt(10.51))
		pdf.Cell(nil, interm.Email)
	}

	// SERVIÇO PRESTADO
	serv := nfse.InfNFSe.DPS.InfDPS.Serv

	pdf.SetFillColor(242, 242, 242)
	pdf.RectFromUpperLeftWithStyle(cmToPt(0.30), cmToPt(12.74), cmToPt(4.90), cmToPt(0.64), "F")
	pdf.SetFillColor(0, 0, 0)

	// LINHA DIVISÓRIA ABAIXO DO INTERMEDIÁRIO (já sobrepõe fundos anteriores e o atual)
	pdf.Line(cmToPt(0.30), cmToPt(12.74), cmToPt(20.70), cmToPt(12.74))

	pdf.SetFont("LiberationSans-Bold", "", 7)
	pdf.SetY(cmToPt(12.81))
	pdf.SetX(cmToPt(0.40))
	pdf.Cell(nil, "SERVIÇO PRESTADO")

	pdf.SetFont("LiberationSans-Bold", "", 6)
	pdf.SetY(cmToPt(12.81))
	pdf.SetX(cmToPt(5.41))
	pdf.Cell(nil, "Código de Tributação Nacional / Municipal")

	pdf.SetX(cmToPt(10.51))
	pdf.Cell(nil, "Código da NBS")

	pdf.SetX(cmToPt(15.62))
	pdf.Cell(nil, "Local da Prestação / Sigla UF / País")

	pdf.SetFont("LiberationSans", "", 7)
	pdf.SetY(cmToPt(13.11))

	tribNac := serv.CServ.CTribNac
	tribMun := serv.CServ.CTribMun
	tribCode := formatTribNac(tribNac)
	if tribMun != "" {
		tribCode += " / " + tribMun
	}
	if tribCode == "" {
		tribCode = "-"
	}
	pdf.SetX(cmToPt(5.41))
	pdf.Cell(nil, tribCode)

	nbs := serv.CServ.CNBS
	if nbs == "" {
		nbs = "-"
	}
	pdf.SetX(cmToPt(10.51))
	pdf.Cell(nil, nbs)

	// O local da prestação na NFS-e autorizada vem do nível infNFSe
	// (<xLocPrestacao>, preenchido pela SEFIN); o <serv><locPrest> da DPS
	// normalmente só traz <cLocPrestacao>.
	loc := nfse.InfNFSe.XLocPrestacao
	if loc == "" {
		loc = serv.LocPrest.XMun
	}
	if loc == "" {
		loc = ibgeCities[serv.LocPrest.CMun]
	}
	if loc == "" {
		loc = serv.LocPrest.CMun
	}
	if loc != "" {
		// Coluna "Local da Prestação / Sigla UF / País": a UF é derivada do código
		// IBGE do local de incidência do ISSQN (a NFS-e não traz a sigla direta).
		if uf := ufFromIBGE(nfse.InfNFSe.CLocIncid); uf != "" {
			loc += " / " + uf
		}
		if serv.LocPrest.CPaisPrestacao != "" {
			loc += " / " + serv.LocPrest.CPaisPrestacao
		}
	} else {
		loc = "-"
	}
	pdf.SetX(cmToPt(15.62))
	pdf.Cell(nil, loc)

	// Descrição do Código de Tributação (Y=13.39 no PDF original)
	pdf.SetFont("LiberationSans", "", 7)
	pdf.SetY(cmToPt(13.46)) // 13.39 + 0.07
	pdf.SetX(cmToPt(0.40))
	descTrib := serv.CServ.XTribMun
	if descTrib == "" {
		descTrib = serv.CServ.XTribNac
	}
	if descTrib == "" {
		// Na NFS-e autorizada, a descrição do código de tributação nacional é
		// preenchida pela SEFIN no nível infNFSe (<xTribNac>), não no <serv>.
		descTrib = nfse.InfNFSe.XTribNac
	}
	if descTrib == "" {
		descTrib = "-"
	}
	pdf.Cell(nil, descTrib)

	// Descrição do Serviço (Y=13.79 no PDF original)
	pdf.SetFont("LiberationSans-Bold", "", 6)
	pdf.SetY(cmToPt(13.86)) // 13.79 + 0.07
	pdf.SetX(cmToPt(0.40))
	pdf.Cell(nil, "Descrição do Serviço")

	pdf.SetFont("LiberationSans", "", 7)
	descServ := serv.CServ.XDescServ
	// Usa SplitText para quebrar em múltiplas linhas, se necessário (largura máx 20cm)
	linhasDesc, errSplit := pdf.SplitText(descServ, cmToPt(20.00))
	if errSplit != nil || len(linhasDesc) == 0 {
		pdf.SetY(cmToPt(14.16))
		pdf.SetX(cmToPt(0.40))
		pdf.Cell(nil, descServ)
	} else {
		// Plota as linhas. O PDF oficial limita o bloco até 14.43, então ~2 linhas cabem bem.
		currentY := 14.16
		for i, linha := range linhasDesc {
			if i >= 2 {
				break
			} // Segurança: limita visualmente a 2 linhas para não invadir o bloco inferior
			pdf.SetY(cmToPt(currentY))
			pdf.SetX(cmToPt(0.40))
			pdf.Cell(nil, linha)
			currentY += 0.35 // Espaçamento de linha simples (~10pt)
		}
	}

	// ----------------------------------------------------------------------
	// TRIBUTAÇÃO MUNICIPAL (ISSQN)
	// ----------------------------------------------------------------------

	tribMunNode := nfse.InfNFSe.DPS.InfDPS.Valores.Trib.TribMun
	issqnApurado := nfse.InfNFSe.Valores.VISSQN

	pdf.SetFillColor(242, 242, 242)
	pdf.RectFromUpperLeftWithStyle(cmToPt(0.30), cmToPt(14.43), cmToPt(4.90), cmToPt(0.64), "F")
	pdf.SetFillColor(0, 0, 0)

	// LINHA DIVISÓRIA ABAIXO DA DESCRIÇÃO DO SERVIÇO (Início do próximo bloco em Y=14.43)
	// Desenhada depois do fundo para sobrepor
	pdf.Line(cmToPt(0.30), cmToPt(14.43), cmToPt(20.70), cmToPt(14.43))

	// Linha 1 (Y=14.43)
	pdf.SetFont("LiberationSans-Bold", "", 7)
	pdf.SetY(cmToPt(14.50))
	pdf.SetX(cmToPt(0.40))
	pdf.Cell(nil, "TRIBUTAÇÃO MUNICIPAL (ISSQN)")

	pdf.SetFont("LiberationSans-Bold", "", 6)
	pdf.SetY(cmToPt(14.50))
	pdf.SetX(cmToPt(5.41))
	pdf.Cell(nil, "Tipo de Tributação do ISSQN")

	pdf.SetX(cmToPt(10.51))
	pdf.Cell(nil, "Município / Sigla UF / País da Incidência do ISSQN")

	pdf.SetFont("LiberationSans", "", 7)
	pdf.SetY(cmToPt(14.80))

	tpTrib := formatTribISSQN(tribMunNode.TribISSQN)
	pdf.SetX(cmToPt(5.41))
	pdf.Cell(nil, tpTrib)

	locIncid := tribMunNode.XLocIncid
	if tribMunNode.CPaisResult != "" {
		locIncid += " / " + tribMunNode.CPaisResult
	}
	if locIncid == "" {
		locIncid = "-"
	}
	pdf.SetX(cmToPt(10.51))
	pdf.Cell(nil, locIncid)

	// Linha 2 (Y=15.08)
	pdf.SetFont("LiberationSans-Bold", "", 6)
	pdf.SetY(cmToPt(15.15))
	pdf.SetX(cmToPt(0.40))
	pdf.Cell(nil, "Regime Especial de Tributação do ISSQN")
	pdf.SetX(cmToPt(5.41))
	pdf.Cell(nil, "Tipo de Imunidade do ISSQN")
	pdf.SetX(cmToPt(10.51))
	pdf.Cell(nil, "Suspensão da Exigibilidade do ISSQN")
	pdf.SetX(cmToPt(15.62))
	pdf.Cell(nil, "Número Processo Suspensão")

	pdf.SetFont("LiberationSans", "", 7)
	pdf.SetY(cmToPt(15.45))
	regEsp := formatRegEspTrib(nfse.InfNFSe.DPS.InfDPS.Prest.RegTrib.RegEspTrib)
	pdf.SetX(cmToPt(0.40))
	pdf.Cell(nil, regEsp)

	imunidade := tribMunNode.TpImunidade
	if imunidade == "" {
		imunidade = "-"
	}
	pdf.SetX(cmToPt(5.41))
	pdf.Cell(nil, imunidade)

	suspensao := formatTpSusp(tribMunNode.ExigSusp.TpSusp)
	pdf.SetX(cmToPt(10.51))
	pdf.Cell(nil, suspensao)

	processo := tribMunNode.ExigSusp.NProcesso
	if processo == "" {
		processo = "-"
	}
	pdf.SetX(cmToPt(15.62))
	pdf.Cell(nil, processo)

	// Linha 3 (Y=15.73)
	pdf.SetFont("LiberationSans-Bold", "", 6)
	pdf.SetY(cmToPt(15.80))
	pdf.SetX(cmToPt(0.40))
	pdf.Cell(nil, "Benefício Municipal")
	pdf.SetX(cmToPt(5.41))
	pdf.Cell(nil, "Cálculo do BM")
	pdf.SetX(cmToPt(10.51))
	pdf.Cell(nil, "Total Deduções/Reduções")
	pdf.SetX(cmToPt(15.62))
	pdf.Cell(nil, "Desconto Incondicionado")

	pdf.SetFont("LiberationSans", "", 7)
	pdf.SetY(cmToPt(16.10))
	bm := nfse.InfNFSe.Valores.TpBM
	if bm == "" {
		bm = "-"
	}
	pdf.SetX(cmToPt(0.40))
	pdf.Cell(nil, bm)

	calcBm := nfse.InfNFSe.Valores.VCalcBM
	if calcBm == "" {
		calcBm = tribMunNode.BM.VRedBCBM
	}
	if calcBm == "" {
		calcBm = "-"
	}
	pdf.SetX(cmToPt(5.41))
	pdf.Cell(nil, calcBm)

	totDed := nfse.InfNFSe.Valores.VCalcDR
	if totDed == "" {
		totDed = nfse.InfNFSe.DPS.InfDPS.Valores.VDedRed.VDR
	}
	if totDed == "" {
		totDed = "-"
	}
	pdf.SetX(cmToPt(10.51))
	pdf.Cell(nil, totDed)

	descInc := nfse.InfNFSe.DPS.InfDPS.Valores.VDescCondIncond.VDescIncond
	if descInc == "" {
		descInc = "-"
	}
	pdf.SetX(cmToPt(15.62))
	pdf.Cell(nil, descInc)

	// Linha 4 (Y=16.37)
	pdf.SetFont("LiberationSans-Bold", "", 6)
	pdf.SetY(cmToPt(16.44))
	pdf.SetX(cmToPt(0.40))
	pdf.Cell(nil, "BC ISSQN")
	pdf.SetX(cmToPt(5.41))
	pdf.Cell(nil, "Alíquota Aplicada")
	pdf.SetX(cmToPt(10.51))
	pdf.Cell(nil, "Retenção do ISSQN")
	pdf.SetX(cmToPt(15.62))
	pdf.Cell(nil, "ISSQN Apurado")

	pdf.SetFont("LiberationSans", "", 7)
	pdf.SetY(cmToPt(16.74))
	bcIssqn := nfse.InfNFSe.Valores.VBC
	if bcIssqn == "" {
		bcIssqn = "-"
	}
	pdf.SetX(cmToPt(0.40))
	pdf.Cell(nil, bcIssqn)

	aliq := nfse.InfNFSe.Valores.PAliqAplic
	if aliq == "" {
		aliq = "-"
	}
	pdf.SetX(cmToPt(5.41))
	pdf.Cell(nil, aliq)

	retIssqn := formatTpRetISSQN(tribMunNode.TpRetISSQN)
	pdf.SetX(cmToPt(10.51))
	pdf.Cell(nil, retIssqn)
	if issqnApurado == "" {
		issqnApurado = "-"
	}
	pdf.SetX(cmToPt(15.62))
	pdf.Cell(nil, issqnApurado)

	// ----------------------------------------------------------------------
	// TRIBUTAÇÃO FEDERAL (EXCETO CBS)
	// ----------------------------------------------------------------------
	tribFed := nfse.InfNFSe.DPS.InfDPS.Valores.Trib.TribFed

	pdf.SetFillColor(242, 242, 242)
	pdf.RectFromUpperLeftWithStyle(cmToPt(0.30), cmToPt(17.02), cmToPt(4.90), cmToPt(0.64), "F")
	pdf.SetFillColor(0, 0, 0)

	// LINHA DIVISÓRIA ABAIXO DO BLOCO ISSQN (Início do próximo bloco em Y=17.02)
	pdf.Line(cmToPt(0.30), cmToPt(17.02), cmToPt(20.70), cmToPt(17.02))

	// Linha 1 (Y=17.02)
	pdf.SetFont("LiberationSans-Bold", "", 7)
	pdf.SetY(cmToPt(17.09))
	pdf.SetX(cmToPt(0.40))
	pdf.Cell(nil, "TRIBUTAÇÃO FEDERAL (EXCETO CBS)")

	pdf.SetFont("LiberationSans-Bold", "", 6)
	pdf.SetX(cmToPt(5.41))
	pdf.Cell(nil, "IRRF")

	pdf.SetX(cmToPt(10.51))
	pdf.Cell(nil, "Contribuição Previdenciária - Retida")

	pdf.SetX(cmToPt(15.62))
	pdf.Cell(nil, "Contribuições Sociais - Retidas")

	pdf.SetFont("LiberationSans", "", 7)
	pdf.SetY(cmToPt(17.39))

	// Mock fallback if empty
	irrf := tribFed.VRetIRRF
	if irrf == "" {
		irrf = "-"
	}
	pdf.SetX(cmToPt(5.41))
	pdf.Cell(nil, irrf)

	cpRet := tribFed.VRetCP
	if cpRet == "" {
		cpRet = "-"
	}
	pdf.SetX(cmToPt(10.51))
	pdf.Cell(nil, cpRet)

	csllRet := tribFed.VRetCSLL
	if csllRet == "" {
		csllRet = "-"
	}
	pdf.SetX(cmToPt(15.62))
	pdf.Cell(nil, csllRet)

	// Linha 2 (Y=17.67)
	pdf.SetFont("LiberationSans-Bold", "", 6)
	pdf.SetY(cmToPt(17.74))
	pdf.SetX(cmToPt(0.40))
	pdf.Cell(nil, "PIS - Débito Apuração Própria")

	pdf.SetX(cmToPt(5.41))
	pdf.Cell(nil, "COFINS - Débito Apuração Própria")

	pdf.SetX(cmToPt(10.51))
	pdf.Cell(nil, "Descrição Contrib. Sociais - Retidas")

	pdf.SetFont("LiberationSans", "", 7)
	pdf.SetY(cmToPt(18.04))

	pis := tribFed.PisCofins.VPis
	if pis == "" {
		pis = "-"
	}
	pdf.SetX(cmToPt(0.40))
	pdf.Cell(nil, pis)

	cofins := tribFed.PisCofins.VCofins
	if cofins == "" {
		cofins = "-"
	}
	pdf.SetX(cmToPt(5.41))
	pdf.Cell(nil, cofins)

	tpret := formatTpRetPisCofins(tribFed.PisCofins.TpRetPisCofins)
	pdf.SetX(cmToPt(10.51))
	pdf.Cell(nil, tpret)

	// ----------------------------------------------------------------------
	// TRIBUTAÇÃO IBS / CBS
	// ----------------------------------------------------------------------
	ibscbs := nfse.InfNFSe.DPS.InfDPS.IBSCBS

	pdf.SetFillColor(242, 242, 242)
	pdf.RectFromUpperLeftWithStyle(cmToPt(0.30), cmToPt(18.32), cmToPt(4.90), cmToPt(0.64), "F")
	pdf.SetFillColor(0, 0, 0)

	// LINHA DIVISÓRIA ABAIXO DA TRIBUTAÇÃO FEDERAL (Y=18.32)
	pdf.Line(cmToPt(0.30), cmToPt(18.32), cmToPt(20.70), cmToPt(18.32))

	// Linha 1 (Y=18.32)
	pdf.SetFont("LiberationSans-Bold", "", 7)
	pdf.SetY(cmToPt(18.39))
	pdf.SetX(cmToPt(0.40))
	pdf.Cell(nil, "TRIBUTAÇÃO IBS / CBS")

	pdf.SetFont("LiberationSans-Bold", "", 6)
	pdf.SetX(cmToPt(5.41))
	pdf.Cell(nil, "CST / CCLASS TRIB")

	pdf.SetX(cmToPt(10.51))
	pdf.Cell(nil, "Indicador Op. / IBGE Incidência / Município Incidência / UF")

	pdf.SetFont("LiberationSans", "", 7)
	pdf.SetY(cmToPt(18.69))

	cstCclass := ibscbs.Valores.Trib.GIBSCBS.CClassTrib
	if cstCclass == "" {
		cstCclass = "-"
	}
	pdf.SetX(cmToPt(5.41))
	pdf.Cell(nil, cstCclass)

	locIbs := ibscbs.CIndOp
	if ibscbs.CLocalidadeIncid != "" {
		locIbs += " / " + ibscbs.CLocalidadeIncid
	}
	if ibscbs.XLocalidadeIncid != "" {
		locIbs += " / " + ibscbs.XLocalidadeIncid
	}
	if locIbs == "" {
		locIbs = "-"
	}
	pdf.SetX(cmToPt(10.51))
	pdf.Cell(nil, locIbs)

	// Linha 2 (Y=18.96)
	pdf.SetFont("LiberationSans-Bold", "", 6)
	pdf.SetY(cmToPt(19.03))
	pdf.SetX(cmToPt(0.40))
	pdf.Cell(nil, "Exclusões e Reduções da Base de Cálculo")

	pdf.SetX(cmToPt(5.41))
	pdf.Cell(nil, "Base de Cálculo Após Exclusões e Reduções")

	pdf.SetX(cmToPt(10.51))
	pdf.Cell(nil, "Red. Alíquota IBS / Red. Alíquota CBS")

	pdf.SetX(cmToPt(15.62))
	pdf.Cell(nil, "Alíquota - IBS UF / IBS MUN")

	pdf.SetFont("LiberationSans", "", 7)
	pdf.SetY(cmToPt(19.33))

	excl := "-" // Somatório
	pdf.SetX(cmToPt(0.40))
	pdf.Cell(nil, excl)

	bcIbscbs := ibscbs.Valores.VBC
	if bcIbscbs == "" {
		bcIbscbs = "-"
	}
	pdf.SetX(cmToPt(5.41))
	pdf.Cell(nil, bcIbscbs)

	redAliq := ibscbs.Valores.UF.PRedAliqUF
	if ibscbs.Valores.Mun.PRedAliqMun != "" {
		redAliq += " / " + ibscbs.Valores.Mun.PRedAliqMun
	}
	if ibscbs.Valores.Fed.PRedAliqCBS != "" {
		redAliq += " / " + ibscbs.Valores.Fed.PRedAliqCBS
	}
	if redAliq == "" {
		redAliq = "-"
	}
	pdf.SetX(cmToPt(10.51))
	pdf.Cell(nil, redAliq)

	aliqIbs := ibscbs.Valores.UF.PIBSUF
	if ibscbs.Valores.Mun.PIBSMun != "" {
		aliqIbs += " / " + ibscbs.Valores.Mun.PIBSMun
	}
	if aliqIbs == "" {
		aliqIbs = "-"
	}
	pdf.SetX(cmToPt(15.62))
	pdf.Cell(nil, aliqIbs)

	// Linha 3 (Y=19.61)
	pdf.SetFont("LiberationSans-Bold", "", 6)
	pdf.SetY(cmToPt(19.68))
	pdf.SetX(cmToPt(0.40))
	pdf.Cell(nil, "Alíq. Efetiva Municipal - IBS")

	pdf.SetX(cmToPt(5.41))
	pdf.Cell(nil, "Valor Apurado Municipal - IBS")

	pdf.SetX(cmToPt(10.51))
	pdf.Cell(nil, "Alíq. Efetiva Estadual - IBS")

	pdf.SetX(cmToPt(15.62))
	pdf.Cell(nil, "Valor Apurado Estadual - IBS")

	pdf.SetFont("LiberationSans", "", 7)
	pdf.SetY(cmToPt(19.98))

	aliqEfetMun := ibscbs.Valores.Mun.PAliqEfetMun
	if aliqEfetMun == "" {
		aliqEfetMun = "-"
	}
	pdf.SetX(cmToPt(0.40))
	pdf.Cell(nil, aliqEfetMun)

	vIbsMun := ibscbs.TotCIBS.GIBS.GIBSMunTot.VIBSMun
	if vIbsMun == "" {
		vIbsMun = "-"
	}
	pdf.SetX(cmToPt(5.41))
	pdf.Cell(nil, vIbsMun)

	aliqEfetUf := ibscbs.Valores.UF.PAliqEfetUF
	if aliqEfetUf == "" {
		aliqEfetUf = "-"
	}
	pdf.SetX(cmToPt(10.51))
	pdf.Cell(nil, aliqEfetUf)

	vIbsUf := ibscbs.TotCIBS.GIBS.GIBSUFTot.VIBSUF
	if vIbsUf == "" {
		vIbsUf = "-"
	}
	pdf.SetX(cmToPt(15.62))
	pdf.Cell(nil, vIbsUf)

	// Linha 4 (Y=20.26)
	pdf.SetFont("LiberationSans-Bold", "", 6)
	pdf.SetY(cmToPt(20.33))
	pdf.SetX(cmToPt(0.40))
	pdf.Cell(nil, "Valor Total Apurado - IBS")

	pdf.SetX(cmToPt(5.41))
	pdf.Cell(nil, "Alíquota - CBS")

	pdf.SetX(cmToPt(10.51))
	pdf.Cell(nil, "Alíquota Efetiva - CBS")

	pdf.SetX(cmToPt(15.62))
	pdf.Cell(nil, "Valor Total Apurado - CBS")

	pdf.SetFont("LiberationSans", "", 7)
	pdf.SetY(cmToPt(20.63))

	vIbsTot := ibscbs.TotCIBS.GIBS.VIBSTot
	if vIbsTot == "" {
		vIbsTot = "-"
	}
	pdf.SetX(cmToPt(0.40))
	pdf.Cell(nil, vIbsTot)

	pCbs := ibscbs.Valores.Fed.PCBS
	if pCbs == "" {
		pCbs = "-"
	}
	pdf.SetX(cmToPt(5.41))
	pdf.Cell(nil, pCbs)

	aliqEfetCbs := ibscbs.Valores.Fed.PAliqEfetCBS
	if aliqEfetCbs == "" {
		aliqEfetCbs = "-"
	}
	pdf.SetX(cmToPt(10.51))
	pdf.Cell(nil, aliqEfetCbs)

	vCbs := ibscbs.TotCIBS.GCBS.VCBS
	if vCbs == "" {
		vCbs = "-"
	}
	pdf.SetX(cmToPt(15.62))
	pdf.Cell(nil, vCbs)

	// ----------------------------------------------------------------------
	// VALOR TOTAL DA NFS-E
	// ----------------------------------------------------------------------
	pdf.SetFillColor(242, 242, 242)
	pdf.RectFromUpperLeftWithStyle(cmToPt(0.30), cmToPt(20.90), cmToPt(4.90), cmToPt(0.67), "F")
	pdf.SetFillColor(0, 0, 0)

	// LINHA DIVISÓRIA FINAL DO BLOCO IBS/CBS (Início de Valor Total) Y=20.90
	pdf.Line(cmToPt(0.30), cmToPt(20.90), cmToPt(20.70), cmToPt(20.90))

	// Linha 1 (Y=20.90)
	pdf.SetFont("LiberationSans-Bold", "", 7)
	pdf.SetY(cmToPt(20.97))
	pdf.SetX(cmToPt(0.40))
	pdf.Cell(nil, "VALOR TOTAL DA NFS-E")

	pdf.SetFont("LiberationSans-Bold", "", 6)
	pdf.SetX(cmToPt(5.41))
	pdf.Cell(nil, "Valor da Operação / Serviço")

	pdf.SetX(cmToPt(10.51))
	pdf.Cell(nil, "Desconto Incondicionado")

	pdf.SetX(cmToPt(15.62))
	pdf.Cell(nil, "Desconto Condicionado")

	pdf.SetFont("LiberationSans", "", 7)
	pdf.SetY(cmToPt(21.27))

	vServ := nfse.InfNFSe.DPS.InfDPS.Valores.VServPrest.VServ
	vServ = formatCurrency(vServ)
	if vServ == "" {
		vServ = "-"
	}
	pdf.SetX(cmToPt(5.41))
	pdf.Cell(nil, vServ)

	vDescIncond := nfse.InfNFSe.DPS.InfDPS.Valores.VDescCondIncond.VDescIncond
	vDescIncond = formatCurrency(vDescIncond)
	if vDescIncond == "" {
		vDescIncond = "-"
	}
	pdf.SetX(cmToPt(10.51))
	pdf.Cell(nil, vDescIncond)

	vDescCond := nfse.InfNFSe.DPS.InfDPS.Valores.VDescCondIncond.VDescCond
	vDescCond = formatCurrency(vDescCond)
	if vDescCond == "" {
		vDescCond = "-"
	}
	pdf.SetX(cmToPt(15.62))
	pdf.Cell(nil, vDescCond)

	// Linha 2 (Y=21.59)
	// Neste bloco não colocamos a linha horizontal (como no modelo original, o título "TOTAL DAS RETENÇÕES" está em X=0.30 mas o fundo cinza só vai até Y=21.57 - NT prevê Y=21.59 para a linha 2, então não desenhamos divisória na linha 2)
	pdf.SetFont("LiberationSans-Bold", "", 6)
	pdf.SetY(cmToPt(21.66))
	pdf.SetX(cmToPt(0.40))
	pdf.Cell(nil, "Total das Retenções (ISSQN / Federais)")

	pdf.SetX(cmToPt(5.41))
	pdf.Cell(nil, "Valor Líquido da NFS-e")

	pdf.SetX(cmToPt(10.51))
	pdf.Cell(nil, "Total do IBS/CBS")

	pdf.SetX(cmToPt(15.62))
	pdf.Cell(nil, "Valor Líquido da NFS-e + IBS/CBS")

	pdf.SetFont("LiberationSans", "", 7)
	pdf.SetY(cmToPt(21.96))

	vTotalRet := nfse.InfNFSe.Valores.VTotalRet
	vTotalRet = formatCurrency(vTotalRet)
	if vTotalRet == "" {
		vTotalRet = "-"
	}
	pdf.SetX(cmToPt(0.40))
	pdf.Cell(nil, vTotalRet)

	vLiq := nfse.InfNFSe.Valores.VLiq
	vLiq = formatCurrency(vLiq)
	if vLiq == "" {
		vLiq = "-"
	}
	pdf.SetX(cmToPt(5.41))
	pdf.Cell(nil, vLiq)

	vTotIbsCbs := "-"
	// TODO: calcular soma se necessário
	pdf.SetX(cmToPt(10.51))
	pdf.Cell(nil, vTotIbsCbs)

	vTotNF := ibscbs.TotCIBS.VTotNF
	vTotNF = formatCurrency(vTotNF)
	if vTotNF == "" {
		vTotNF = "-"
	}
	pdf.SetX(cmToPt(15.62))
	pdf.Cell(nil, vTotNF)

	// LINHA DIVISÓRIA FINAL DO BLOCO (Y=22.27)
	pdf.Line(cmToPt(0.30), cmToPt(22.27), cmToPt(20.70), cmToPt(22.27))

	// ----------------------------------------------------------------------
	// INFORMAÇÕES COMPLEMENTARES
	// ----------------------------------------------------------------------
	pdf.SetFont("LiberationSans-Bold", "", 7)
	pdf.SetY(cmToPt(22.34))
	pdf.SetX(cmToPt(0.40))
	pdf.Cell(nil, "INFORMAÇÕES COMPLEMENTARES")

	// Conteúdo Y=22.68
	pdf.SetFont("LiberationSans", "", 7)
	pdf.SetY(cmToPt(22.68))
	xOutInf := nfse.InfNFSe.DPS.InfDPS.Serv.InfoCompl.XOutInf
	// Se houver texto, plota
	if xOutInf != "" {
		pdf.SetX(cmToPt(0.40))
		pdf.MultiCell(&gopdf.Rect{W: cmToPt(20.00), H: cmToPt(4.00)}, xOutInf)
	}

	// ----------------------------------------------------------------------
	// DATA CIENTIFICAÇÃO E ASSINATURA
	// ----------------------------------------------------------------------
	pdf.Line(cmToPt(0.30), cmToPt(28.10), cmToPt(20.70), cmToPt(28.10))  // Linha superior
	pdf.Line(cmToPt(0.30), cmToPt(28.77), cmToPt(20.70), cmToPt(28.77))  // Linha inferior
	pdf.Line(cmToPt(0.30), cmToPt(28.10), cmToPt(0.30), cmToPt(28.77))   // Borda esq
	pdf.Line(cmToPt(20.70), cmToPt(28.10), cmToPt(20.70), cmToPt(28.77)) // Borda dir

	// Divisórias verticais
	pdf.Line(cmToPt(5.41), cmToPt(28.10), cmToPt(5.41), cmToPt(28.77))
	pdf.Line(cmToPt(10.51), cmToPt(28.10), cmToPt(10.51), cmToPt(28.77))

	pdf.SetFont("LiberationSans-Bold", "", 6)
	pdf.SetY(cmToPt(28.17))
	pdf.SetX(cmToPt(0.40))
	pdf.Cell(nil, "DATA CIENTIFICAÇÃO:")

	pdf.SetX(cmToPt(5.51))
	pdf.Cell(nil, "IDENTIFICAÇÃO E ASSINATURA")

	pdf.SetX(cmToPt(10.61))
	pdf.Cell(nil, "Nº NFS-e / CHAVE NFS-e")

	pdf.SetFont("LiberationSans", "", 6)
	pdf.SetY(cmToPt(28.47))
	nNFSe := nfse.InfNFSe.NNFSe
	chave := strings.TrimPrefix(nfse.InfNFSe.ID, "NFS")
	if nNFSe == "" {
		nNFSe = "-"
	}
	if chave == "" {
		chave = "-"
	}

	pdf.SetX(cmToPt(10.61))
	pdf.Cell(nil, nNFSe+" / "+chave)

	var buf bytes.Buffer
	_, err = pdf.WriteTo(&buf)
	if err != nil {
		return nil, fmt.Errorf("pdf: erro ao renderizar: %w", err)
	}
	return buf.Bytes(), nil
}
