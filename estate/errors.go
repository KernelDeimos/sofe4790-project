package estate

import "errors"

type ErrorState struct {
	Last    error
	History []error
}

func (state *ErrorState) Add(err error) {
	if err != nil {
		state.Last = err
		state.History = append(state.History, err)
	}
}

func (state *ErrorState) AddNewError(msg string) {
	state.Add(errors.New(msg))
}

func (state *ErrorState) Error() string {
	switch len(state.History) {
	case 0:
		return ""
	case 1:
		return state.Last.Error()
	default:
		text := "error list:\n"
		for _, err := range state.History {
			text += "\terror item: " + err.Error() + "\n"
		}
		return text
	}
}

func (state *ErrorState) GetError() error {
	switch len(state.History) {
	case 0:
		return nil
	default:
		return state
	}
}
