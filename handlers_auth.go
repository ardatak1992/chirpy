package main

import (
	"net/http"
	"time"

	"github.com/ardatak1992/chirpy/internal/auth"
	"github.com/ardatak1992/chirpy/internal/database"
)

func (cfg *apiConfig) handlerUserLogin(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type response struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	params := parameters{}
	err := decodeRequestJSON(r, &params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error while logging in", err)
		return
	}

	expirationTime := time.Hour

	userRes, err := cfg.dbQueries.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error while logging in", err)
		return
	}

	token, err := auth.MakeJWT(userRes.ID, cfg.tokenSecret, expirationTime)
	if err != nil {
		respondWithError(w, http.StatusOK, "error while logging", err)
	}

	passwordCorrect, err := auth.CheckPasswordHash(params.Password, userRes.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error while logging in", err)
		return
	}

	if !passwordCorrect {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", nil)
		return
	}

	refreshTokenStr := auth.MakeRefreshToken()

	refreshToken, err := cfg.dbQueries.CreateRefreshToken(
		r.Context(),
		database.CreateRefreshTokenParams{
			Token:  refreshTokenStr,
			UserID: userRes.ID,
		},
	)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error while logging in", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:        userRes.ID,
			CreatedAt: userRes.CreatedAt,
			UpdatedAt: userRes.UpdatedAt,
			Email:     userRes.Email,
		},
		Token:        token,
		RefreshToken: refreshToken.Token,
	})

}

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	refreshTokenStr, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error refreshing", err)
	}

	refreshToken, err := cfg.dbQueries.GetRefreshToken(r.Context(), refreshTokenStr)

}
