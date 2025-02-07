package qstnnr

import "github.com/mateopresacastro/qstnnr/pkg/store"

func getInitialData() store.InitialData {
	questions := map[store.QuestionID]store.Question{
		1: {
			ID:   1,
			Text: "What function is used for deferred execution in Go?",
			Options: map[store.OptionID]store.Option{
				1: {ID: 1, Text: "wait()"},
				2: {ID: 2, Text: "defer()"},
				3: {ID: 3, Text: "delayed()"},
				4: {ID: 4, Text: "async()"},
			},
		},
		2: {
			ID:   2,
			Text: "Which of these is the correct way to declare a slice in Go?",
			Options: map[store.OptionID]store.Option{
				1: {ID: 1, Text: "var s array[]int"},
				2: {ID: 2, Text: "var s []int"},
				3: {ID: 3, Text: "s := array{int}"},
				4: {ID: 4, Text: "s := list[int]"},
			},
		},
		3: {
			ID:   3,
			Text: "What is the zero value for a pointer in Go?",
			Options: map[store.OptionID]store.Option{
				1: {ID: 1, Text: "nil"},
				2: {ID: 2, Text: "0"},
				3: {ID: 3, Text: "undefined"},
				4: {ID: 4, Text: "void"},
			},
		},
		4: {
			ID:   4,
			Text: "Which keyword is used to create a new goroutine?",
			Options: map[store.OptionID]store.Option{
				1: {ID: 1, Text: "go"},
				2: {ID: 2, Text: "goroutine"},
				3: {ID: 3, Text: "routine"},
				4: {ID: 4, Text: "async"},
			},
		},
		5: {
			ID:   5,
			Text: "What happens if you try to send to a closed channel in Go?",
			Options: map[store.OptionID]store.Option{
				1: {ID: 1, Text: "The program will panic"},
				2: {ID: 2, Text: "The send will block"},
				3: {ID: 3, Text: "The value is discarded silently"},
				4: {ID: 4, Text: "A runtime error occurs without panic"},
			},
		},
		6: {
			ID:   6,
			Text: "Which of these correctly declares a variable that can hold any type in Go?",
			Options: map[store.OptionID]store.Option{
				1: {ID: 1, Text: "var x interface{}"},
				2: {ID: 2, Text: "var x any"},
				3: {ID: 3, Text: "Both both interface{} and any are correct"},
				4: {ID: 4, Text: "var x object"},
			},
		},
		7: {
			ID:   7,
			Text: "What is the purpose of the blank identifier (_) in Go?",
			Options: map[store.OptionID]store.Option{
				1: {ID: 1, Text: "To discard an unwanted value"},
				2: {ID: 2, Text: "To declare a private variable"},
				3: {ID: 3, Text: "To create an anonymous function"},
				4: {ID: 4, Text: "To mark a variable as nullable"},
			},
		},
		8: {
			ID:   8,
			Text: "How do you make a field in a struct unexported in Go?",
			Options: map[store.OptionID]store.Option{
				1: {ID: 1, Text: "Start the field name with a lowercase letter"},
				2: {ID: 2, Text: "Use the private keyword"},
				3: {ID: 3, Text: "Add an underscore prefix"},
				4: {ID: 4, Text: "Add the unexported tag"},
			},
		},
		9: {
			ID:   9,
			Text: "What is the correct way to check if a key exists in a map?",
			Options: map[store.OptionID]store.Option{
				1: {ID: 1, Text: "value, exists := map[key]"},
				2: {ID: 2, Text: "exists := key in map"},
				3: {ID: 3, Text: "exists := map.contains(key)"},
				4: {ID: 4, Text: "exists := map.has(key)"},
			},
		},
		10: {
			ID:   10,
			Text: "Which of these correctly implements an empty interface?",
			Options: map[store.OptionID]store.Option{
				1: {ID: 1, Text: "type I interface {}"},
				2: {ID: 2, Text: "type I interface { void }"},
				3: {ID: 3, Text: "type I = interface"},
				4: {ID: 4, Text: "interface I {}"},
			},
		},
	}

	solutions := map[store.QuestionID]store.OptionID{
		1:  2, // defer()
		2:  2, // var s []int
		3:  1, // nil
		4:  1, // go
		5:  1, // panic
		6:  3, // both interface{} and any are correct
		7:  1, // discard unwanted value
		8:  1, // lowercase letter
		9:  1, // value, exists := map[key]
		10: 1, // type I interface {}
	}

	return store.InitialData{
		Questions: questions,
		Solutions: solutions,
	}
}
