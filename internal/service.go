package service

import (
	"net/http"

	"github.com/hashicorp/go-hclog"
	"github.com/labstack/echo/v4"
	"github.com/manedurphy/boggle-solver/pkg/boggle"
)

type (
	Config struct {
		Logger hclog.Logger
	}

	Service struct {
		cfg Config
	}

	SolveRequest struct {
		Board [][]string `json:"board"`
	}

	SolveResponse struct {
		Result *boggle.Result `json:"result"`
	}

	ErrorMessage struct {
		Message string `json:"message"`
	}
)

func New(cfg Config) *Service {
	return &Service{
		cfg: cfg,
	}
}

func (s *Service) Solve(c echo.Context) error {
	var (
		b      *boggle.Boggle
		result *boggle.Result
		req    SolveRequest
		err    error
	)

	s.cfg.Logger.With("func", "Solve").Info("incoming request")

	if err = c.Bind(&req); err != nil {
		return c.JSON(http.StatusInternalServerError, &ErrorMessage{
			Message: "internal server error",
		})
	}
	s.cfg.Logger.Debug("boggle board", "board", req.Board)

	b, err = boggle.New(boggle.Config{
		Board:        req.Board,
		WordsZipFile: "/home/dane/Documents/interview_prep/assessments/rstudio/data/words.zip",
	})
	if err != nil {
		s.cfg.Logger.With("err", err, "board", req.Board).Error("could not create new boggle game")
		return c.JSON(http.StatusBadRequest, &ErrorMessage{
			Message: "invalid boggle board submitted",
		})
	}

	result = b.Solve()
	s.cfg.Logger.Info("successfully solved boggle board!")

	return c.JSON(http.StatusOK, &SolveResponse{
		Result: result,
	})
}
