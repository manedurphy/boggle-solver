package service

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/manedurphy/boggle-solver/pkg/boggle"
)

type (
	Service struct {
	}

	SolveRequest struct {
		Board [][]string `json:"board"`
	}

	SolveResponse struct {
		WordsFound []string `json:"words_found"`
	}

	ErrorMessage struct {
		Message string `json:"message"`
	}
)

func New() *Service {
	return &Service{}
}

func (s *Service) Solve(c echo.Context) error {
	var (
		b          *boggle.Boggle
		req        SolveRequest
		wordsFound []string
		err        error
	)

	fmt.Println("incoming request")

	if err = c.Bind(&req); err != nil {
		return c.JSON(http.StatusInternalServerError, &ErrorMessage{
			Message: "internal server error",
		})
	}

	b, err = boggle.New(req.Board)
	if err != nil {
		fmt.Printf("could not create new boggle game: %s\n", err)
		return c.JSON(http.StatusBadRequest, &ErrorMessage{
			Message: "invalid boggle board submitted",
		})
	}

	wordsFound = b.Solve()
	return c.JSON(http.StatusOK, &SolveResponse{
		WordsFound: wordsFound,
	})
}
