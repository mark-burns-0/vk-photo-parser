package parser

type Configer interface {
	GetOwnerID() string
	GetAlbumID() string
	GetToken() string
}

func New(config Configer) *Parser {
	return &Parser{}
}

func (p *Parser) Parse() {

}
