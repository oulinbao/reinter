// This program is free software: you can redistribute it and/or modify it
// under the terms of the GNU General Public License as published by the Free
// Software Foundation, either version 3 of the License, or (at your option)
// any later version.
//
// This program is distributed in the hope that it will be useful, but
// WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General
// Public License for more details.
//
// You should have received a copy of the GNU General Public License along
// with this program.  If not, see <http://www.gnu.org/licenses/>.

package intersection

import (
	"testing"
)

func TestIntersection(t *testing.T) {
	type Case struct {
		Expr1  string
		Expr2  string
		Expect bool
	}
	cases := []Case{
		{"a*bba+", "b*aaab+a", false},
		{"a*ba+", "b*ab+a", true}, //aba
		{"/api/v1/[0-9]+/get", `/api/v1/\w+/get`, true},
		{"/api/v1/[0-9]+", `/api/v1/[a-z]+`, false},
		{"/api/v1/[0-9]+/", `/api/v1/\w+/`, true},
		{"/api/v1/[0-9]+/get", `/api/v1/[a-zA-Z]+/get`, false},
		{"/api/v1/[0-9]+/get", `/api/v1/[0-9a-z]+/get`, true},
		{"/api/v1/[0-9]+/get", `/api/v1/\w+/get`, true},
		{" ", `\s`, true},
		{"[a-zA-Z]+", "[a-z]+", true},
		{"[a-zA-Z]+", "[a-z]+", true},
		{"(a|b)", "(b|c)", true},     // Positive case: intersection 'b'
		{"(a|b|c)", "(c|d|e)", true}, // Positive case: intersection 'c'
		{"a*bba+", "b*aaabbb+a", false},
		{"", "", true},
		{"a+", "a?", true},
		{"a*", "a+", true},
		{"a*", "a?", true},

		{"a*", "b*", true},

		//
		{"[A-Z]+", "[a-z]+", false},
		{"a", "b", false},

		{"\\s+", "a+", false},
		{"/api/v1/.*/", "/api/v2/.*/", false},
		{"/api/v1/[0-9]+/get", "/api/v1/[a-z]+/get", false},

		{"api/v1/a/b", "api/v1/[a-z]/b", true},
		{"webapi/v1/session/[\\d]{1,3}/items/", "webapi/v1/session/1/items/.*", true},
		{"", "", true},
		{"/api/v1/[0-9]+/get", "/api/v1/[a-z]+/get", false},
		{"a*b", "c*d", false},
		{"^abc$", "^abc$", true}, // Positive case: exact match
		{"a+", "a?", true},
		{"a*", "a+", true},
		{"a*", "a?", true},
		{"a*", "b*", true},

		//
		{"[A-Z]+", "[a-z]+", false},
		{"a", "b", false},
		{"\\s+", "a+", false},
		{"/api/v1/.*/", "/api/v2/.*/", false},
		{"", "a", false}, // Negative case: no intersection

		{"a*b", "ab*", true}, // Positive case: both regex match 'ab'
		{"a*", "a", true},    // Positive case: 'a*' includes 'a'
		// Negative case: empty regex and non-empty
		{"a", "", false},           // Negative case: non-empty regex and empty
		{"(a|b)", "(c|d)", false},  // Negative case: no intersection
		{"(a|b)*", "(b|c)*", true}, // Positive case: intersection possible
		{"a+", "b+", false},        // Negative case: no intersection
		{"[0-9]", "[0-9]", true},   // Positive case: same regex
		{"[0-9]", "[a-z]", false},  // Negative case: no intersection
		{"a", "a*", true},          // Positive case: 'a' is included in 'a*'
		{"^abc$", "^def$", false},  // Negative case: no match
		{"abc", "abc", true},       // Positive case: exact match without anchors
		{"abc", "xyz", false},      // Negative case: no match
		{"(a|b)", "(d|e)", false},  // Negative case: no
	}

	for _, tt := range cases {
		t.Run(tt.Expr1+" ∩ "+tt.Expr2, func(t *testing.T) {
			result, _ := HasIntersection(tt.Expr1, tt.Expr2)

			if result != tt.Expect {
				t.Errorf("expected %v, got %v", tt.Expect, result)
			}
		})
	}

}
