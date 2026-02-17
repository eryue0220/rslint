package no_dupe_keys

import (
	"github.com/microsoft/typescript-go/shim/ast"
	"github.com/web-infra-dev/rslint/internal/rule"
	"github.com/web-infra-dev/rslint/internal/utils"
)

// Message builder
func buildDupeKeysMessage(name string) rule.RuleMessage {
	return rule.RuleMessage{
		Id:          "unexpected",
		Description: "Duplicate key '" + name + "'.",
	}
}

// seenFlags tracks what property types have been seen for each key.
// Getter + setter for the same key is allowed; all other duplicates are errors.
const (
	seenData = 1 << iota
	seenGet
	seenSet
)

func shouldReport(property *ast.Node, seen map[string]int, nameText string) bool {
	prev := seen[nameText]
	report := false

	switch property.Kind {
	case ast.KindGetAccessor:
		if prev != 0 && prev != seenSet {
			report = true
		}
		seen[nameText] = prev | seenGet
	case ast.KindSetAccessor:
		if prev != 0 && prev != seenGet {
			report = true
		}
		seen[nameText] = prev | seenSet
	default:
		if prev != 0 {
			report = true
		}
		seen[nameText] = prev | seenData
	}

	return report
}

func normalizeKey(sourceFile *ast.SourceFile, name *ast.Node) (key, nameText string, nameType utils.MemberNameType) {
	nameText, nameType = utils.GetNameFromMember(sourceFile, name)
	key = nameText

	// Remove surrounding quotes from string literals
	if nameType == utils.MemberNameTypeQuoted && len(nameText) >= 2 && nameText[0] == '"' && nameText[len(nameText)-1] == '"' {
		key = nameText[1 : len(nameText)-1]
	} else if name.Kind == ast.KindBigIntLiteral && len(nameText) > 0 && nameText[len(nameText)-1] == 'n' {
		// Remove 'n' from BigInt literals
		key = nameText[:len(nameText)-1]
	}

	return key, nameText, nameType
}

var NoDupeKeysRule = rule.CreateRule(rule.Rule{
	Name: "no-dupe-keys",
	Run: func(ctx rule.RuleContext, options any) rule.RuleListeners {
		return rule.RuleListeners{
			ast.KindObjectLiteralExpression: func(node *ast.Node) {
				objectLiteral := node.AsObjectLiteralExpression()
				if objectLiteral == nil {
					return
				}

				length := len(objectLiteral.Properties.Nodes)
				if length < 2 {
					return
				}

				seen := make(map[string]int)

				for _, property := range objectLiteral.Properties.Nodes {
					name := property.Name()
					if name == nil {
						continue
					}

					key, nameText, nameType := normalizeKey(ctx.SourceFile, name)
					if key == "__proto__" && name.Kind != ast.KindComputedPropertyName && property.Kind == ast.KindPropertyAssignment {
						continue
					}

					if nameType == utils.MemberNameTypeExpression {
						continue
					}

					if shouldReport(property, seen, key) {
						ctx.ReportNode(property, buildDupeKeysMessage(nameText))
					}
				}
			},
		}
	},
})
