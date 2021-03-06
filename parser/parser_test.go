package parser_test

import (
	"testing"

	"github.com/BytecodeAgency/import-boundary-checker/lexer"
	"github.com/BytecodeAgency/import-boundary-checker/parser"
	"github.com/stretchr/testify/assert"
)

func TestParser_Parse_Correct(t *testing.T) {
	tests := []struct {
		input              string
		expectedLang       parser.Language
		expectedImportBase string
		expectedRules      []parser.Rule
	}{
		{`LANG "Go";

IMPORTRULE "github.com/BytecodeAgency/someexampleproject/platform-backend/typings/entities"
CANNOTIMPORT "github.com/BytecodeAgency/someexampleproject/platform-backend";

IMPORTRULE
  	"github.com/BytecodeAgency/someexampleproject/platform-backend/domain"
CANNOTIMPORT
	"github.com/BytecodeAgency/someexampleproject/platform-backend/infrastructure"
	"github.com/BytecodeAgency/someexampleproject/platform-backend/data"
ALLOW
	"github.com/BytecodeAgency/someexampleproject/platform-backend/infrastructure/detail";`,
			"Go",
			"",
			[]parser.Rule{
				{"github.com/BytecodeAgency/someexampleproject/platform-backend/typings/entities",
					[]string{"github.com/BytecodeAgency/someexampleproject/platform-backend"},
					[]string{}},
				{"github.com/BytecodeAgency/someexampleproject/platform-backend/domain",
					[]string{
						"github.com/BytecodeAgency/someexampleproject/platform-backend/infrastructure",
						"github.com/BytecodeAgency/someexampleproject/platform-backend/data"},
					[]string{
						"github.com/BytecodeAgency/someexampleproject/platform-backend/infrastructure/detail"}},
			}},
		{`LANG "Go";
IMPORTBASE "github.com/BytecodeAgency/someexampleproject/platform-backend";

IMPORTRULE "github.com/BytecodeAgency/someexampleproject/platform-backend/typings/entities"
CANNOTIMPORT "github.com/BytecodeAgency/someexampleproject/platform-backend"
ALLOW
	"github.com/BytecodeAgency/someexampleproject/platform-backend/platform-backend/detail";

IMPORTRULE
  	"github.com/BytecodeAgency/someexampleproject/platform-backend/domain"
CANNOTIMPORT
	"github.com/BytecodeAgency/someexampleproject/platform-backend/infrastructure"
	"github.com/BytecodeAgency/someexampleproject/platform-backend/data"
ALLOW
	"github.com/BytecodeAgency/someexampleproject/platform-backend/platform-backend/detail"
	"github.com/BytecodeAgency/someexampleproject/platform-backend/infrastructure/detail";`,
			"Go",
			"github.com/BytecodeAgency/someexampleproject/platform-backend",
			[]parser.Rule{
				{"github.com/BytecodeAgency/someexampleproject/platform-backend/typings/entities",
					[]string{"github.com/BytecodeAgency/someexampleproject/platform-backend"},
					[]string{
						"github.com/BytecodeAgency/someexampleproject/platform-backend/platform-backend/detail",
					}},
				{"github.com/BytecodeAgency/someexampleproject/platform-backend/domain",
					[]string{
						"github.com/BytecodeAgency/someexampleproject/platform-backend/infrastructure",
						"github.com/BytecodeAgency/someexampleproject/platform-backend/data"},
					[]string{
						"github.com/BytecodeAgency/someexampleproject/platform-backend/platform-backend/detail",
						"github.com/BytecodeAgency/someexampleproject/platform-backend/infrastructure/detail",
					}},
			}},
		{`LANG "Go";
IMPORTBASE "github.com/BytecodeAgency/someexampleproject/platform-backend";

IMPORTRULE "[IMPORTBASE]/typings/entities"
CANNOTIMPORT "[IMPORTBASE]"
ALLOW "[IMPORTBASE]/data";

IMPORTRULE "[IMPORTBASE]/domain"
CANNOTIMPORT
	"[IMPORTBASE]/infrastructure"
	"[IMPORTBASE]/data"
ALLOW
	"[IMPORTBASE]/infrastructure/detail"
	"[IMPORTBASE]/data/detail";`,
			"Go",
			"github.com/BytecodeAgency/someexampleproject/platform-backend",
			[]parser.Rule{
				{"github.com/BytecodeAgency/someexampleproject/platform-backend/typings/entities",
					[]string{"github.com/BytecodeAgency/someexampleproject/platform-backend"},
					[]string{"github.com/BytecodeAgency/someexampleproject/platform-backend/data"}},
				{"github.com/BytecodeAgency/someexampleproject/platform-backend/domain",
					[]string{
						"github.com/BytecodeAgency/someexampleproject/platform-backend/infrastructure",
						"github.com/BytecodeAgency/someexampleproject/platform-backend/data"},
					[]string{
						"github.com/BytecodeAgency/someexampleproject/platform-backend/infrastructure/detail",
						"github.com/BytecodeAgency/someexampleproject/platform-backend/data/detail"}},
			}},
	}

	for _, test := range tests {
		// Lexer
		l := lexer.New(test.input)
		l.Exec()
		res, errs := l.Result()
		assert.Empty(t, errs)

		// Parser
		p := parser.New(res)
		p.Parse()
		assert.Empty(t, p.Errors)
		assert.Equal(t, test.expectedLang, p.Lang)
		assert.Equal(t, test.expectedRules, p.Rules)
		assert.Equal(t, test.expectedImportBase, p.ImportBase)
	}
}

func TestParser_Parse_Incorrect(t *testing.T) {
	incorrectInputs := []string{
		// Invalid language
		`LANG "COBOL";
IMPORTRULE "some/module"
CANNOTIMPORT "some/other/module";`,

		// Multiple importrules
		`LANG "Go";
IMPORTRULE "some/module1" "some/module2"
CANNOTIMPORT "some/other/module";`,

		// Not finishing the importrule
		`LANG "Go";
IMPORTRULE "some/module1" "some/module2";`,

		// Not setting the importrule, only the cannotimports
		`LANG "Go";
IMPORTRULE
CANNOTIMPORT "some/module2";`,

		// Not setting the language
		`IMPORTRULE "some/module"
CANNOTIMPORT "some/module2";`,

		// Only setting the language, and no importrules
		`LANG "Go";`,
	}

	for _, input := range incorrectInputs {
		// Lexer
		l := lexer.New(input)
		l.Exec()
		res, errs := l.Result()
		assert.Empty(t, errs)

		// Parser
		p := parser.New(res)
		p.Parse()
		assert.NotEmpty(t, p.Errors)
	}
}
