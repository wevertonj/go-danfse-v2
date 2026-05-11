package danfse

import (
	"encoding/xml"
	"fmt"
)

type NFSe struct {
	XMLName xml.Name    `xml:"http://www.sped.fazenda.gov.br/nfse NFSe"`
	InfNFSe InfNFSeNode `xml:"http://www.sped.fazenda.gov.br/nfse infNFSe"`
}

type InfNFSeNode struct {
	ID            string   `xml:"Id,attr"`
	XLocEmi       string   `xml:"http://www.sped.fazenda.gov.br/nfse xLocEmi"`
	NNFSe         string   `xml:"http://www.sped.fazenda.gov.br/nfse nNFSe"`
	CLocIncid     string   `xml:"http://www.sped.fazenda.gov.br/nfse cLocIncid"`
	XTribNac      string   `xml:"http://www.sped.fazenda.gov.br/nfse xTribNac"`
	AmbGer        string   `xml:"http://www.sped.fazenda.gov.br/nfse ambGer"`
	DhProc        string   `xml:"http://www.sped.fazenda.gov.br/nfse dhProc"`
	Emit          EmitNode `xml:"http://www.sped.fazenda.gov.br/nfse emit"`
	Valores       Valores  `xml:"http://www.sped.fazenda.gov.br/nfse valores"`
	DPS           DPSNode  `xml:"http://www.sped.fazenda.gov.br/nfse DPS"`
}

type EmitNode struct {
	CNPJ     string       `xml:"http://www.sped.fazenda.gov.br/nfse CNPJ"`
	XNome    string       `xml:"http://www.sped.fazenda.gov.br/nfse xNome"`
	EnderNac EnderNacNode `xml:"http://www.sped.fazenda.gov.br/nfse enderNac"`
	Fone     string       `xml:"http://www.sped.fazenda.gov.br/nfse fone"`
	Email    string       `xml:"http://www.sped.fazenda.gov.br/nfse email"`
}

type EnderNacNode struct {
	XLgr    string `xml:"http://www.sped.fazenda.gov.br/nfse xLgr"`
	Nro     string `xml:"http://www.sped.fazenda.gov.br/nfse nro"`
	XCpl    string `xml:"http://www.sped.fazenda.gov.br/nfse xCpl"`
	XBairro string `xml:"http://www.sped.fazenda.gov.br/nfse xBairro"`
	CMun    string `xml:"http://www.sped.fazenda.gov.br/nfse cMun"`
	XMun    string `xml:"http://www.sped.fazenda.gov.br/nfse xMun"`
	UF      string `xml:"http://www.sped.fazenda.gov.br/nfse UF"`
	CEP     string `xml:"http://www.sped.fazenda.gov.br/nfse CEP"`
}

type Valores struct {
	VLiq       string `xml:"http://www.sped.fazenda.gov.br/nfse vLiq"`
	VTotalRet  string `xml:"http://www.sped.fazenda.gov.br/nfse vTotalRet"`
	TpBM       string `xml:"http://www.sped.fazenda.gov.br/nfse tpBM"`
	VCalcBM    string `xml:"http://www.sped.fazenda.gov.br/nfse vCalcBM"`
	VCalcDR    string `xml:"http://www.sped.fazenda.gov.br/nfse vCalcDR"`
	VBC        string `xml:"http://www.sped.fazenda.gov.br/nfse vBC"`
	PAliqAplic string `xml:"http://www.sped.fazenda.gov.br/nfse pAliqAplic"`
	VISSQN     string `xml:"http://www.sped.fazenda.gov.br/nfse vISSQN"`
}

type DPSNode struct {
	InfDPS InfDPSNode `xml:"http://www.sped.fazenda.gov.br/nfse infDPS"`
}

type InfDPSNode struct {
	TpAmb   string     `xml:"http://www.sped.fazenda.gov.br/nfse tpAmb"`
	DhEmi   string     `xml:"http://www.sped.fazenda.gov.br/nfse dhEmi"`
	Serie   string     `xml:"http://www.sped.fazenda.gov.br/nfse serie"`
	NDPS    string     `xml:"http://www.sped.fazenda.gov.br/nfse nDPS"`
	DCompet string     `xml:"http://www.sped.fazenda.gov.br/nfse dCompet"`
	Prest   PrestNode  `xml:"http://www.sped.fazenda.gov.br/nfse prest"`
	Toma    TomaNode   `xml:"http://www.sped.fazenda.gov.br/nfse toma"`
	Interm  IntermNode `xml:"http://www.sped.fazenda.gov.br/nfse interm"`
	Serv    ServNode   `xml:"http://www.sped.fazenda.gov.br/nfse serv"`
	IBSCBS  IBSCBSNode `xml:"http://www.sped.fazenda.gov.br/nfse IBSCBS"`
	Valores ValoresDPS `xml:"http://www.sped.fazenda.gov.br/nfse valores"`
}

type PrestNode struct {
	CNPJ  string `xml:"http://www.sped.fazenda.gov.br/nfse CNPJ"`
	IM    string `xml:"http://www.sped.fazenda.gov.br/nfse IM"`
	Fone  string `xml:"http://www.sped.fazenda.gov.br/nfse fone"`
	Email string `xml:"http://www.sped.fazenda.gov.br/nfse email"`
	RegTrib struct {
		OpSimpNac   string `xml:"http://www.sped.fazenda.gov.br/nfse opSimpNac"`
		RegApTribSN string `xml:"http://www.sped.fazenda.gov.br/nfse regApTribSN"`
		RegEspTrib  string `xml:"http://www.sped.fazenda.gov.br/nfse regEspTrib"`
	} `xml:"http://www.sped.fazenda.gov.br/nfse regTrib"`
}

type TomaNode struct {
	CNPJ  string `xml:"http://www.sped.fazenda.gov.br/nfse CNPJ"`
	XNome string `xml:"http://www.sped.fazenda.gov.br/nfse xNome"`
	IM    string `xml:"http://www.sped.fazenda.gov.br/nfse IM"`
	Fone  string `xml:"http://www.sped.fazenda.gov.br/nfse fone"`
	Email string `xml:"http://www.sped.fazenda.gov.br/nfse email"`
	End   struct {
		EndNac EnderNacNode `xml:"http://www.sped.fazenda.gov.br/nfse endNac"`
		XLgr   string       `xml:"http://www.sped.fazenda.gov.br/nfse xLgr"`
		Nro    string       `xml:"http://www.sped.fazenda.gov.br/nfse nro"`
		XCpl   string       `xml:"http://www.sped.fazenda.gov.br/nfse xCpl"`
		XBairro string      `xml:"http://www.sped.fazenda.gov.br/nfse xBairro"`
	} `xml:"http://www.sped.fazenda.gov.br/nfse end"`
}

type IntermNode struct {
	CNPJ  string `xml:"http://www.sped.fazenda.gov.br/nfse CNPJ"`
	XNome string `xml:"http://www.sped.fazenda.gov.br/nfse xNome"`
	IM    string `xml:"http://www.sped.fazenda.gov.br/nfse IM"`
	Fone  string `xml:"http://www.sped.fazenda.gov.br/nfse fone"`
	Email string `xml:"http://www.sped.fazenda.gov.br/nfse email"`
	CMun  string `xml:"http://www.sped.fazenda.gov.br/nfse cMun"`
	XMun  string `xml:"http://www.sped.fazenda.gov.br/nfse xMun"`
	CEP   string `xml:"http://www.sped.fazenda.gov.br/nfse CEP"`
	UF    string `xml:"http://www.sped.fazenda.gov.br/nfse UF"`
	End   struct {
		EndNac EnderNacNode `xml:"http://www.sped.fazenda.gov.br/nfse endNac"`
		XLgr   string       `xml:"http://www.sped.fazenda.gov.br/nfse xLgr"`
		Nro    string       `xml:"http://www.sped.fazenda.gov.br/nfse nro"`
		XCpl   string       `xml:"http://www.sped.fazenda.gov.br/nfse xCpl"`
		XBairro string      `xml:"http://www.sped.fazenda.gov.br/nfse xBairro"`
	} `xml:"http://www.sped.fazenda.gov.br/nfse end"`
}

type ServNode struct {
	LocPrest struct {
		XLocPrestacao  string `xml:"http://www.sped.fazenda.gov.br/nfse xLocPrestacao"`
		CPaisPrestacao string `xml:"http://www.sped.fazenda.gov.br/nfse cPaisPrestacao"`
		CMun           string `xml:"http://www.sped.fazenda.gov.br/nfse cMun"`
		XMun           string `xml:"http://www.sped.fazenda.gov.br/nfse xMun"`
		CEP            string `xml:"http://www.sped.fazenda.gov.br/nfse CEP"`
	} `xml:"http://www.sped.fazenda.gov.br/nfse locPrest"`
	Obra struct {
		CObra string `xml:"http://www.sped.fazenda.gov.br/nfse cObra"`
		Art   string `xml:"http://www.sped.fazenda.gov.br/nfse art"`
	} `xml:"http://www.sped.fazenda.gov.br/nfse obra"`
	CServ struct {
		CTribNac  string `xml:"http://www.sped.fazenda.gov.br/nfse cTribNac"`
		CTribMun  string `xml:"http://www.sped.fazenda.gov.br/nfse cTribMun"`
		CNBS      string `xml:"http://www.sped.fazenda.gov.br/nfse cNBS"`
		XTribNac  string `xml:"http://www.sped.fazenda.gov.br/nfse xTribNac"`
		XTribMun  string `xml:"http://www.sped.fazenda.gov.br/nfse xTribMun"`
		XDescServ string `xml:"http://www.sped.fazenda.gov.br/nfse xDescServ"`
	} `xml:"http://www.sped.fazenda.gov.br/nfse cServ"`
	InfoCompl InfoComplNode `xml:"http://www.sped.fazenda.gov.br/nfse infoCompl"`
}

type InfoComplNode struct {
	XOutInf string `xml:"http://www.sped.fazenda.gov.br/nfse xOutInf"`
}

type DestNode struct {
	CNPJ  string `xml:"http://www.sped.fazenda.gov.br/nfse CNPJ"`
	XNome string `xml:"http://www.sped.fazenda.gov.br/nfse xNome"`
	Fone  string `xml:"http://www.sped.fazenda.gov.br/nfse fone"`
	Email string `xml:"http://www.sped.fazenda.gov.br/nfse email"`
	CMun  string `xml:"http://www.sped.fazenda.gov.br/nfse cMun"`
	XMun  string `xml:"http://www.sped.fazenda.gov.br/nfse xMun"`
	CEP   string `xml:"http://www.sped.fazenda.gov.br/nfse CEP"`
	UF    string `xml:"http://www.sped.fazenda.gov.br/nfse UF"`
	End   struct {
		EndNac EnderNacNode `xml:"http://www.sped.fazenda.gov.br/nfse endNac"`
		XLgr   string       `xml:"http://www.sped.fazenda.gov.br/nfse xLgr"`
		Nro    string       `xml:"http://www.sped.fazenda.gov.br/nfse nro"`
		XCpl   string       `xml:"http://www.sped.fazenda.gov.br/nfse xCpl"`
		XBairro string      `xml:"http://www.sped.fazenda.gov.br/nfse xBairro"`
	} `xml:"http://www.sped.fazenda.gov.br/nfse end"`
}

type IBSCBSNode struct {
	Dest DestNode `xml:"http://www.sped.fazenda.gov.br/nfse dest"`
	CIndOp           string `xml:"http://www.sped.fazenda.gov.br/nfse cIndOp"`
	CLocalidadeIncid string `xml:"http://www.sped.fazenda.gov.br/nfse cLocalidadeIncid"`
	XLocalidadeIncid string `xml:"http://www.sped.fazenda.gov.br/nfse xLocalidadeIncid"`
	Valores          struct {
		VDescIncond string `xml:"http://www.sped.fazenda.gov.br/nfse vDescIncond"`
		VCalcReeRepRes string `xml:"http://www.sped.fazenda.gov.br/nfse vCalcReeRepRes"`
		VBC string `xml:"http://www.sped.fazenda.gov.br/nfse vBC"`
		UF struct {
			PRedAliqUF string `xml:"http://www.sped.fazenda.gov.br/nfse pRedAliqUF"`
			PIBSUF string `xml:"http://www.sped.fazenda.gov.br/nfse pIBSUF"`
			PAliqEfetUF string `xml:"http://www.sped.fazenda.gov.br/nfse pAliqEfetUF"`
		} `xml:"http://www.sped.fazenda.gov.br/nfse uf"`
		Mun struct {
			PRedAliqMun string `xml:"http://www.sped.fazenda.gov.br/nfse pRedAliqMun"`
			PIBSMun string `xml:"http://www.sped.fazenda.gov.br/nfse pIBSMun"`
			PAliqEfetMun string `xml:"http://www.sped.fazenda.gov.br/nfse pAliqEfetMun"`
		} `xml:"http://www.sped.fazenda.gov.br/nfse mun"`
		Fed struct {
			PRedAliqCBS string `xml:"http://www.sped.fazenda.gov.br/nfse pRedAliqCBS"`
			PCBS string `xml:"http://www.sped.fazenda.gov.br/nfse pCBS"`
			PAliqEfetCBS string `xml:"http://www.sped.fazenda.gov.br/nfse pAliqEfetCBS"`
		} `xml:"http://www.sped.fazenda.gov.br/nfse fed"`
		Trib struct {
			GIBSCBS struct {
				CClassTrib string `xml:"http://www.sped.fazenda.gov.br/nfse cClassTrib"`
			} `xml:"http://www.sped.fazenda.gov.br/nfse gIBSCBS"`
		} `xml:"http://www.sped.fazenda.gov.br/nfse trib"`
	} `xml:"http://www.sped.fazenda.gov.br/nfse valores"`
	TotCIBS struct {
		GIBS struct {
			GIBSMunTot struct {
				VIBSMun string `xml:"http://www.sped.fazenda.gov.br/nfse vIBSMun"`
			} `xml:"http://www.sped.fazenda.gov.br/nfse gIBSMunTot"`
			GIBSUFTot struct {
				VIBSUF string `xml:"http://www.sped.fazenda.gov.br/nfse vIBSUF"`
			} `xml:"http://www.sped.fazenda.gov.br/nfse gIBSUFTot"`
			VIBSTot string `xml:"http://www.sped.fazenda.gov.br/nfse vIBSTot"`
		} `xml:"http://www.sped.fazenda.gov.br/nfse gIBS"`
		GCBS struct {
			VCBS string `xml:"http://www.sped.fazenda.gov.br/nfse vCBS"`
		} `xml:"http://www.sped.fazenda.gov.br/nfse gCBS"`
		VTotNF string `xml:"http://www.sped.fazenda.gov.br/nfse vTotNF"`
	} `xml:"http://www.sped.fazenda.gov.br/nfse totCIBS"`
}

type ValoresDPS struct {
	VServPrest struct {
		VServ string `xml:"http://www.sped.fazenda.gov.br/nfse vServ"`
	} `xml:"http://www.sped.fazenda.gov.br/nfse vServPrest"`
	Trib struct {
		TotTrib struct {
			PTotTribSN string `xml:"http://www.sped.fazenda.gov.br/nfse pTotTribSN"`
		} `xml:"http://www.sped.fazenda.gov.br/nfse totTrib"`
		TribMun struct {
			TribISSQN   string `xml:"http://www.sped.fazenda.gov.br/nfse tribISSQN"`
			CLocIncid   string `xml:"http://www.sped.fazenda.gov.br/nfse cLocIncid"`
			XLocIncid   string `xml:"http://www.sped.fazenda.gov.br/nfse xLocIncid"`
			CPaisResult string `xml:"http://www.sped.fazenda.gov.br/nfse cPaisResult"`
			TpImunidade string `xml:"http://www.sped.fazenda.gov.br/nfse tpImunidade"`
			ExigSusp struct {
				TpSusp    string `xml:"http://www.sped.fazenda.gov.br/nfse tpSusp"`
				NProcesso string `xml:"http://www.sped.fazenda.gov.br/nfse nProcesso"`
			} `xml:"http://www.sped.fazenda.gov.br/nfse exigSusp"`
			TpRetISSQN string `xml:"http://www.sped.fazenda.gov.br/nfse tpRetISSQN"`
			BM struct {
				VRedBCBM string `xml:"http://www.sped.fazenda.gov.br/nfse vRedBCBM"`
			} `xml:"http://www.sped.fazenda.gov.br/nfse BM"`
		} `xml:"http://www.sped.fazenda.gov.br/nfse tribMun"`
		TribFed struct {
			VRetIRRF string `xml:"http://www.sped.fazenda.gov.br/nfse vRetIRRF"`
			VRetCP   string `xml:"http://www.sped.fazenda.gov.br/nfse vRetCP"`
			VRetCSLL string `xml:"http://www.sped.fazenda.gov.br/nfse vRetCSLL"`
			PisCofins struct {
				VPis           string `xml:"http://www.sped.fazenda.gov.br/nfse vPis"`
				VCofins        string `xml:"http://www.sped.fazenda.gov.br/nfse vCofins"`
				TpRetPisCofins string `xml:"http://www.sped.fazenda.gov.br/nfse tpRetPisCofins"`
			} `xml:"http://www.sped.fazenda.gov.br/nfse piscofins"`
		} `xml:"http://www.sped.fazenda.gov.br/nfse tribFed"`
	} `xml:"http://www.sped.fazenda.gov.br/nfse trib"`
	VDedRed struct {
		VDR string `xml:"http://www.sped.fazenda.gov.br/nfse vDR"`
	} `xml:"http://www.sped.fazenda.gov.br/nfse vDedRed"`
	VDescCondIncond struct {
		VDescIncond string `xml:"http://www.sped.fazenda.gov.br/nfse vDescIncond"`
		VDescCond   string `xml:"http://www.sped.fazenda.gov.br/nfse vDescCond"`
	} `xml:"http://www.sped.fazenda.gov.br/nfse vDescCondIncond"`
}

// ParseXML converte o XML da NFS-e em structs para uso na renderização
func ParseXML(xmlData []byte) (*NFSe, error) {
	var n NFSe
	if err := xml.Unmarshal(xmlData, &n); err != nil {
		return nil, fmt.Errorf("falha ao parsear xml: %w", err)
	}
	return &n, nil
}
