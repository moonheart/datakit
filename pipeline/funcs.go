package pipeline

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"os"

	"github.com/tidwall/gjson"

	"gitlab.jiagouyun.com/cloudcare-tools/datakit/pipeline/parser"
)

const (
	CONTENT = "__content"
)

type ProFunc func(p *Pipeline, node parser.Node) (*Pipeline, error)

var (
	funcsMap = map[string]ProFunc{
		"grok":      Grok,
		"json":      Json,
		"rename":    Rename,
		"strfmt":    Strfmt,
		"cast":      Cast,
		"expr":      Expr,

		"user_agent": UserAgent,
		"url_decode": UrlDecode,
		"geoip":     GeoIp,
		"datetime":  DateTime,
		"group":     Group,
		"group_in":  GroupIn,
	}
)

func Json(p *Pipeline, node parser.Node) (*Pipeline, error) {
	funcExpr := node.(*parser.FuncExpr)
	if len(funcExpr.Param) < 2 || len(funcExpr.Param) > 3{
		return nil, fmt.Errorf("func %s expected 2 or 3 args", funcExpr.Name)
	}

	key := funcExpr.Param[0].(*parser.Identifier).Name
	old := funcExpr.Param[1].(*parser.Identifier).Name

	newkey := old
	if len(funcExpr.Param) == 3 {
		newkey = funcExpr.Param[2].(*parser.Identifier).Name
	}

	cont := p.getContentStr(key)
	v := getGjsonResult(cont, old)
	p.setContent(newkey, v)

	return p, nil
}

func Rename(p *Pipeline, node parser.Node) (*Pipeline, error) {
	funcExpr := node.(*parser.FuncExpr)
	if len(funcExpr.Param) != 2 {
		return p, fmt.Errorf("func %s expected 2 args", funcExpr.Name)
	}

	old := funcExpr.Param[1].(*parser.Identifier).Name
	new := funcExpr.Param[0].(*parser.Identifier).Name

	v := p.getContent(old)
	p.setContent(new, v)
	delete(p.Output, old)

	return p, nil
}

func UserAgent(p *Pipeline, node parser.Node) (*Pipeline, error) {
	funcExpr := node.(*parser.FuncExpr)
	if len(funcExpr.Param) != 1 {
		return nil, fmt.Errorf("func %s expected 1 args", funcExpr.Name)
	}

	key := funcExpr.Param[0].(*parser.Identifier).Name

	v := UserAgentHandle(p.getContentStr(key))
	p.setContent(new, v)

	return p, nil
}

func UrlDecode(p *Pipeline, node parser.Node) (*Pipeline, error) {
	funcExpr := node.(*parser.FuncExpr)
	if len(funcExpr.Param) != 1 {
		return nil, fmt.Errorf("func %s expected 1 args", funcExpr.Name)
	}

	key := funcExpr.Param[0].(*parser.Identifier).Name

	if v, err := UrldecodeHandle(p.getContentStr(key)); err != nil {
		return p, err
	} else {
		p.setContent(key, v)
	}

	return p, nil
}

func GeoIp(p *Pipeline, node parser.Node) (*Pipeline, error) {
	funcExpr := node.(*parser.FuncExpr)
	if len(funcExpr.Param) != 1 {
		return nil, fmt.Errorf("func %s expected 2 args", funcExpr.Name)
	}

	key := funcExpr.Param[0].(*parser.Identifier).Name

	if v, err := GeoIpHandle(p.getContentStr(key)); err != nil {
		return p, err
	} else {
		p.setContent(key, v)
	}

	return p, nil
}

func DateTime(p *Pipeline, node parser.Node) (*Pipeline, error) {
	if len(funcExpr.Param) != 3 || len(funcExpr.Param) != 4 {
		return nil, fmt.Errorf("func %s expected 3 or 4 args", funcExpr.Name)
	}

	var tz = 8

	funcExpr := node.(*parser.FuncExpr)

	key       := funcExpr.Param[0].(*parser.Identifier).Name
	precision := funcExpr.Param[1].(*parser.StringLiteral).Val
	fmts      := funcExpr.Param[2].(*parser.StringLiteral).Val

	if len(funcExpr.Param) == 4 {
		tzStr := funcExpr.Param[3].(*parser.StringLiteral).Val
		if v, ok := tzStr.(*parser.NumberLiteral); ok {
			if v.IsInt {
				tz = v.Int
			}
		}
	}

	if v, err := DateFormatHandle(p.getContent(key), 0, fmts, tz); err != nil {
		return p, err
	} else {
		p.setContent(key, v)
	}

	return p, nil
}

func Expr(p *Pipeline, node parser.Node) (*Pipeline, error) {
	funcExpr := node.(*parser.FuncExpr)
	if len(funcExpr.Param) != 2 {
		return nil, fmt.Errorf("func %s expected 2 args", funcExpr.Name)
	}

	key  := funcExpr.Param[1].(*parser.Identifier).Name
	expr := funcExpr.Param[0].(*parser.BinaryExpr)

	if v, err := Calc(expr, p); err != nil {
		return p, err
	} else {
		p.setContent(key, v)
	}


	return p, nil
}

func Strfmt(p *Pipeline, node parser.Node) (*Pipeline, error) {
	outdata := make([]interface{}, 0)

	funcExpr := node.(*parser.FuncExpr)
	if len(funcExpr.Param) < 2 {
		return nil, fmt.Errorf("func %s expected more than 2 args", funcExpr.Name)
	}

	key  := funcExpr.Param[0].(*parser.Identifier).Name
	fmts := funcExpr.Param[1].(*parser.StringLiteral).Val


	for i := 2; i < len(funcExpr.Param); i++ {
		switch v := funcExpr.Param[i].(type) {
		case *parser.Identifier:
			outdata = append(outdata, p.getContent(v.Name))
		case *parser.NumberLiteral:
			if v.IsInt {
				outdata = append(outdata, v.Int)
			} else {
				outdata = append(outdata, v.Float)
			}
		case *parser.StringLiteral:
			outdata = append(outdata, v.Val)
		default:
			outdata = append(outdata, v)
		}
	}

	strfmt := fmt.Sprintf(fmts, outdata...)
    p.setContent(key, strfmt)

	return p, nil
}

func Cast(p *Pipeline, node parser.Node) (*Pipeline, error) {
	funcExpr := node.(*parser.FuncExpr)
	if len(funcExpr.Param) != 2 {
		return nil, fmt.Errorf("func %s expected 2 args", funcExpr.Name)
	}

	key      := funcExpr.Param[0].(*parser.Identifier).Name
	castType := funcExpr.Param[1].(*parser.StringLiteral).Val

	v := cast(p.getContent(key), castType)
	p.setContent(key, v)

	return p, nil
}

func Group(p *Pipeline, node parser.Node) (*Pipeline, error) {
	funcExpr := node.(*parser.FuncExpr)
	if len(funcExpr.Param) != 3 || len(funcExpr.Param) != 4 {
		return nil, fmt.Errorf("func %s expected 3 or 4 args", funcExpr.Name)
	}

	var nvalue interface{}

	funcExpr := node.(*parser.FuncExpr)

	key   := funcExpr.Param[0].(*parser.Identifier).Name
	set   := funcExpr.Param[1].(*parser.NodeList)
	value := funcExpr.Param[2]

	newkey := key
	var start, end float64

	if len(funcExpr.Param) == 4 {
		newkey = funcExpr.Param[3].(*parser.Identifier).Name
	}

	if len(set) != 2 {
		nil, fmt.Errorf("range value is not expected 3 or 4 args", set)
	}

	if v. ok := set[0].(*parser.NumberLiteral); ok {
		if v.IsInt {
			start = float64(v.IsInt)
		} else {
			start = v.Float
		}
	}

	if v. ok := set[1].(*parser.NumberLiteral); ok {
		if v.IsInt {
			end = float64(v.IsInt)
		} else {
			end = v.Float
		}
	}

	if GroupHandle(p.getContent(key), start, end) {
		switch v := funcExpr.Param[i].(type) {
		case *parser.NumberLiteral:
			if v.IsInt {
				p.setContent(newkey, v.IsInt)
			} else {
				p.setContent(newkey, v.Float)
			}
		case *parser.StringLiteral:
			p.setContent(newkey, v.Val)
		case *parser.BoolLiteral:
			p.setContent(newkey, v.Val)
		}
	}

	return p, nil
}

func GroupIn(p *Pipeline, node parser.Node) (*Pipeline, error) {
	setdata := make([]interface{}, 0)
	funcExpr := node.(*parser.FuncExpr)
	if len(funcExpr.Param) != 3 || len(funcExpr.Param) != 4 {
		return nil, fmt.Errorf("func %s expected 3 or 4 args", funcExpr.Name)
	}

	var nvalue interface{}

	funcExpr := node.(*parser.FuncExpr)

	key   := funcExpr.Param[0].(*parser.Identifier).Name
	set   := funcExpr.Param[1].(*parser.NodeList)
	value := funcExpr.Param[2]

	newkey := key
	if len(funcExpr.Param) == 4 {
		newkey = funcExpr.Param[3].(*parser.Identifier).Name
	}

	if len(set) != 2 {
		nil, fmt.Errorf("range value is not expected 3 or 4 args", set)
	}

	for node. ok := range set {
		switch v := node.(type) {
		case *parser.Identifier:
			setdata = append(setdata, p.getContent(v.Name))
		case *parser.NumberLiteral:
			if v.IsInt {
				setdata = append(setdata, v.Int)
			} else {
				setdata = append(setdata, v.Float)
			}
		case *parser.StringLiteral:
			setdata = append(setdata, v.Val)
		default:
			setdata = append(setdata, v)
		}
	}

	if v. ok := set[1].(*parser.NumberLiteral); ok {
		if v.IsInt {
			end = float64(v.IsInt)
		} else {
			end = v.Float
		}
	}

	if GroupHandle(p.getContent(key), setdata, value) {
		switch v := funcExpr.Param[i].(type) {
		case *parser.NumberLiteral:
			if v.IsInt {
				p.setContent(newkey, v.IsInt)
			} else {
				p.setContent(newkey, v.Float)
			}
		case *parser.StringLiteral:
			p.setContent(newkey, v.Val)
		case *parser.BoolLiteral:
			p.setContent(newkey, v.Val)
		}
	}

	return p, nil
}

func ParseScript(scriptOrPath string) ([]parser.Node, error) {
	data := scriptOrPath

	_, err := os.Stat(scriptOrPath)
	if err ==  nil || !os.IsNotExist(err){
		cont, err := ioutil.ReadFile(scriptOrPath)
		if err != nil {
			return nil, err
		}
		data = string(cont)
	}

	nodes, err := parser.ParseFuncExpr(string(data))
	for _, node := range nodes {
		switch v := node.(type) {
		case *parser.FuncExpr:
			DebugNodesHelp(v, "")
		default:
		}
	}

	return nodes, err
}

func DebugNodesHelp(f *parser.FuncExpr, prev string) {
	l.Debugf("%v%v", prev, f.Name)

	for _, node := range f.Param {
		switch v := node.(type) {
		case *parser.FuncExpr:
			DebugNodesHelp(v, prev+"    ")
		default:
			l.Debugf("%v%v", prev+"    |", node)
		}
	}
}

func getGjsonResult(data, id string) interface{} {
	g := gjson.Get(data, id)
	switch g.Type {
	case gjson.Null:
		return nil

	case gjson.False:
		return false

	case gjson.Number:
		if strings.Contains(g.Raw, ".") {
			return g.Float()
		} else {
			return g.Int()
		}

	case gjson.String:
		return g.String()

	case gjson.True:
		return true

	case gjson.JSON:
		return g.Raw

	default:
		return nil
	}
}
