package schema

import (
	. "github.com/parvez3019/go-swagger3/openApi3Schema"
	"github.com/parvez3019/go-swagger3/parser/model"
	"github.com/parvez3019/go-swagger3/parser/utils"
	"go/ast"
	goParser "go/parser"
	"go/token"
	"os"
	"strings"
)

// parseGenericTypeName detects a generic instantiation like "response.ResponseData[model.CreditResponse]".
// Returns base type, type arguments, and true when the type is generic.
func parseGenericTypeName(typeName string) (string, []string, bool) {
	idx := strings.Index(typeName, "[")
	if idx <= 0 {
		return "", nil, false
	}
	lastIdx := strings.LastIndex(typeName, "]")
	if lastIdx != len(typeName)-1 {
		return "", nil, false
	}
	baseType := typeName[:idx]
	typeArgsStr := typeName[idx+1 : lastIdx]
	if typeArgsStr == "" {
		return "", nil, false // empty brackets are handled as arrays elsewhere
	}
	typeArgs := splitTypeArgs(typeArgsStr)
	return baseType, typeArgs, true
}

// splitTypeArgs splits comma-separated type arguments, respecting nested brackets.
func splitTypeArgs(s string) []string {
	var args []string
	depth := 0
	start := 0
	for i, c := range s {
		switch c {
		case '[':
			depth++
		case ']':
			depth--
		case ',':
			if depth == 0 {
				args = append(args, strings.TrimSpace(s[start:i]))
				start = i + 1
			}
		}
	}
	if start < len(s) {
		args = append(args, strings.TrimSpace(s[start:]))
	}
	return args
}

type Parser interface {
	GetPkgAst(pkgPath string) (map[string]*ast.Package, error)
	RegisterType(pkgPath, pkgName, typeName string) (string, error)
	ParseSchemaObject(pkgPath, pkgName, typeName string) (*SchemaObject, error)
}

type parser struct {
	model.Utils
	OpenAPI *OpenAPIObject
}

func NewParser(utils model.Utils, openAPIObject *OpenAPIObject) Parser {
	return &parser{
		Utils:   utils,
		OpenAPI: openAPIObject,
	}
}

func (p *parser) GetPkgAst(pkgPath string) (map[string]*ast.Package, error) {
	if cache, ok := p.PkgPathAstPkgCache[pkgPath]; ok {
		return cache, nil
	}
	ignoreFileFilter := func(info os.FileInfo) bool {
		name := info.Name()
		return !info.IsDir() && !strings.HasPrefix(name, ".") && strings.HasSuffix(name, ".go") && !strings.HasSuffix(name, "_test.go")
	}
	astPackages, err := goParser.ParseDir(token.NewFileSet(), pkgPath, ignoreFileFilter, goParser.ParseComments)
	if err != nil {
		return nil, err
	}
	p.PkgPathAstPkgCache[pkgPath] = astPackages
	return astPackages, nil
}

func (p *parser) RegisterType(pkgPath, pkgName, typeName string) (string, error) {
	var registerTypeName string

	if utils.IsBasicGoType(typeName) || utils.IsInterfaceType(typeName) {
		registerTypeName = typeName
	} else if schemaObject, ok := p.KnownIDSchema[utils.GenSchemaObjectID(pkgName, typeName, p.SchemaWithoutPkg)]; ok {
		_, ok := p.OpenAPI.Components.Schemas[utils.ReplaceBackslash(typeName)]
		if !ok {
			p.OpenAPI.Components.Schemas[utils.ReplaceBackslash(typeName)] = schemaObject
		}
		return utils.GenSchemaObjectID(pkgName, typeName, p.SchemaWithoutPkg), nil
	} else {
		schemaObject, err := p.ParseSchemaObject(pkgPath, pkgName, typeName)
		if err != nil {
			return "", err
		}
		registerTypeName = schemaObject.ID
		_, ok := p.OpenAPI.Components.Schemas[utils.ReplaceBackslash(registerTypeName)]
		if !ok {
			p.OpenAPI.Components.Schemas[utils.ReplaceBackslash(registerTypeName)] = schemaObject
		}
	}
	return registerTypeName, nil
}

func (p *parser) ParseSchemaObject(pkgPath, pkgName, typeName string) (*SchemaObject, error) {
	schemaObject, err, isBasicType := p.parseBasicTypeSchemaObject(pkgPath, pkgName, typeName)
	if isBasicType {
		return schemaObject, err
	}

	// Handle generic type instantiation like "response.ResponseData[model.CreditResponse]"
	if baseType, typeArgs, ok := parseGenericTypeName(typeName); ok {
		return p.parseGenericTypeSchemaObject(pkgPath, pkgName, baseType, typeArgs)
	}

	return p.parseCustomTypeSchemaObject(pkgPath, pkgName, typeName)
}
