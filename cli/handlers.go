package cli

import "fmt"

func handlerLogin(s *state, cmd Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("'login' expects a single argument")
	}
	if err := s.cfg.SetUser(cmd.Args[0]); err != nil {
		return err
	}
	fmt.Println("the user has been set.")
	return nil
}
