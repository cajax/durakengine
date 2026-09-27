package game

import (
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"unicode"

	"github.com/cajax/durakengine/pkg/game"
)

func TestRuleErrorsWorkWithErrorsIsAndAs(t *testing.T) {
	g, players := newTestGame(nil,
		[]*game.Card{card(game.Six, game.Spades), card(game.Seven, game.Clubs)},
		[]*game.Card{card(game.Eight, game.Spades)},
	)
	mustSucceed(t, g.Attack(players[0], []*game.Card{card(game.Six, game.Spades)}))

	err := g.Attack(players[0], []*game.Card{card(game.Seven, game.Clubs)})
	if !errors.Is(err, game.ErrAttackRankNotOnTable) {
		t.Fatalf("Expected errors.Is(err, ErrAttackRankNotOnTable), got %v", err)
	}
	if errors.Is(err, game.ErrAttackIsTooBig) {
		t.Error("errors.Is matched an unrelated rule error")
	}
	var ruleErr *game.RuleError
	if !errors.As(err, &ruleErr) {
		t.Fatalf("Expected errors.As to find a *RuleError in %T", err)
	}
	if ruleErr.Code() != "ATTACK_RANK_NOT_ON_TABLE" {
		t.Errorf("Expected code ATTACK_RANK_NOT_ON_TABLE, got %q", ruleErr.Code())
	}
	if err.Error() != game.ErrorAttackRankNotOnTable {
		t.Errorf("Expected message %q, got %q", game.ErrorAttackRankNotOnTable, err.Error())
	}
}

func TestLifecycleErrorsAreRuleErrors(t *testing.T) {
	g, _ := newTestGame(nil, []*game.Card{}, []*game.Card{})
	if err := g.StartGame(); !errors.Is(err, game.ErrGameAlreadyStarted) {
		t.Errorf("StartGame: expected ErrGameAlreadyStarted, got %v", err)
	}
	if _, err := g.AddPlayer(game.NewPlayer("3", false, false, false, "Player", nil, false, false)); !errors.Is(err, game.ErrAddPlayerWhileInGame) || err.Error() != "cannot_add_player_while_in_game" {
		t.Errorf("AddPlayer: expected ErrAddPlayerWhileInGame, got %v", err)
	}
	if _, err := (&game.Deck{}).GetTrump(); !errors.Is(err, game.ErrDeckNoTrump) {
		t.Errorf("GetTrump: expected ErrDeckNoTrump, got %v", err)
	}
}

// TestEveryErrorConstantHasRuleError checks the engine source: each Error* message constant must have an
// Err* sentinel built from it, with the constant's name without "Error" in UPPER_SNAKE as its code, and no
// error may be returned with errors.New or fmt.Errorf.
func TestEveryErrorConstantHasRuleError(t *testing.T) {
	dir := filepath.Join("..", "..", "pkg", "game")
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	fset := token.NewFileSet()
	constants := map[string]bool{}
	sentinels := map[string][2]string{} // var name -> {code, message constant}
	for _, e := range entries {
		if !strings.HasSuffix(e.Name(), ".go") || strings.HasSuffix(e.Name(), "_test.go") {
			continue
		}
		f, err := parser.ParseFile(fset, filepath.Join(dir, e.Name()), nil, 0)
		if err != nil {
			t.Fatal(err)
		}
		ast.Inspect(f, func(n ast.Node) bool {
			switch n := n.(type) {
			case *ast.GenDecl:
				for _, spec := range n.Specs {
					vs, ok := spec.(*ast.ValueSpec)
					if !ok {
						continue
					}
					for i, name := range vs.Names {
						if n.Tok == token.CONST && strings.HasPrefix(name.Name, "Error") {
							constants[name.Name] = true
						}
						if n.Tok == token.VAR && strings.HasPrefix(name.Name, "Err") && i < len(vs.Values) {
							call, ok := vs.Values[i].(*ast.CallExpr)
							if !ok || len(call.Args) != 2 {
								continue
							}
							code, _ := call.Args[0].(*ast.BasicLit)
							msg, _ := call.Args[1].(*ast.Ident)
							if code == nil || msg == nil {
								continue
							}
							unquoted, _ := strconv.Unquote(code.Value)
							sentinels[name.Name] = [2]string{unquoted, msg.Name}
						}
					}
				}
			case *ast.SelectorExpr:
				if pkg, ok := n.X.(*ast.Ident); ok && (pkg.Name == "errors" && n.Sel.Name == "New" || pkg.Name == "fmt" && n.Sel.Name == "Errorf") {
					t.Errorf("%s: untyped error created with %s.%s", fset.Position(n.Pos()), pkg.Name, n.Sel.Name)
				}
			}
			return true
		})
	}
	if len(constants) == 0 {
		t.Fatal("No Error* constants found")
	}
	for constant := range constants {
		name := strings.TrimPrefix(constant, "Error")
		sentinel, ok := sentinels["Err"+name]
		if !ok {
			t.Errorf("%s has no Err%s sentinel", constant, name)
			continue
		}
		if sentinel[1] != constant {
			t.Errorf("Err%s uses message %s, want %s", name, sentinel[1], constant)
		}
		if want := upperSnake(name); sentinel[0] != want {
			t.Errorf("Err%s has code %q, want %q", name, sentinel[0], want)
		}
	}
}

func upperSnake(s string) string {
	var b strings.Builder
	for i, r := range s {
		if i > 0 && unicode.IsUpper(r) {
			b.WriteByte('_')
		}
		b.WriteRune(unicode.ToUpper(r))
	}
	return b.String()
}
