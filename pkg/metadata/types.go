package metadata

const (
	InactiveCertStr = "TCL"
	ActiveCertStr   = "TCI"
)

type CertificateInfo struct {
	Status       CertStatusType
	GoldStandard string
	Weight       string
	CertNumber   string
	CustomerID   string
	IssueDate    string
	ActiveStatus string
}

type CertStatusType uint32

const (
	CertStatusType_INACTIVE CertStatusType = 0
	CertStatusType_ACTIVE   CertStatusType = 1
)
