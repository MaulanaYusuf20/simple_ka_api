package main

import (
	"math/rand"
	"net/http"
	"os"
	"simple_ka_api/connection"
	"strconv"

	"github.com/labstack/echo"
)

func main() {
	e := echo.New()

	db, err := connection.GetConnection()
	if err != nil {
		panic(err)
	}

	e.POST("/login", func(ctx echo.Context) error {
		type LoginRequest struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		type LoginResponse struct {
			Role string `json:"role"`
		}

		type Response struct {
			Code   int            `json:"code"`
			Status string         `json:"status"`
			Data   *LoginResponse `json:"data"`
		}

		// get body request
		req := new(LoginRequest)
		err := ctx.Bind(req)

		if err != nil {
			return ctx.JSON(http.StatusBadRequest, &Response{
				Code:   http.StatusBadRequest,
				Status: err.Error(),
			})
		}

		// query
		SQL := "SELECT role FROM public.login WHERE login.email = ($1) AND login.password = ($2)"
		rows, err := db.Query(SQL, req.Email, req.Password)
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, &Response{
				Code:   http.StatusBadRequest,
				Status: err.Error(),
			})
		}

		if rows.Next() {
			var role string
			rows.Scan(&role)

			return ctx.JSON(http.StatusOK, &Response{
				Code:   http.StatusOK,
				Status: "success",
				Data: &LoginResponse{
					Role: role,
				},
			})
		}

		return ctx.JSON(http.StatusUnauthorized, &Response{
			Code:   http.StatusUnauthorized,
			Status: "success",
		})
	})

	e.POST("/register", func(ctx echo.Context) error {
		type RegisterRequest struct {
			Nik      string `json:"nik"`
			Nama     string `json:"nama"`
			NoHp     string `json:"noHp"`
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		type Response struct {
			Code   int    `json:"code"`
			Status string `json:"status"`
		}

		// get body request
		req := new(RegisterRequest)
		err := ctx.Bind(req)
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, &Response{
				Code:   http.StatusBadRequest,
				Status: err.Error(),
			})
		}

		// add data to customer
		SQL := "INSERT INTO public.customer (nik, nama, email, no_hp) VALUES ($1, $2, $3, $4)"
		_, err = db.Exec(SQL, req.Nik, req.Nama, req.Email, req.NoHp)
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, &Response{
				Code:   http.StatusBadRequest,
				Status: err.Error(),
			})
		}

		SQL2 := "INSERT INTO public.login (email, password) VALUES ($1, $2)"
		_, err = db.Exec(SQL2, req.Email, req.Password)
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, &Response{
				Code:   http.StatusBadRequest,
				Status: err.Error(),
			})
		}

		return ctx.JSON(http.StatusOK, &Response{
			Code:   http.StatusOK,
			Status: "success add data",
		})
	})

	e.GET("/jadwal", func(ctx echo.Context) error {
		type JadwalModel struct {
			IdKereta   string `json:"idKereta"`
			HargaTiket int    `json:"hargaTiket"`
			NamaKereta string `json:"namaKereta"`
			Tujuan     string `json:"tujuan"`
			Asal       string `json:"asal"`
		}

		type JadwalResponse struct {
			Code   int            `json:"code"`
			Status string         `json:"status"`
			Data   []*JadwalModel `json:"data"`
		}

		SQL := "SELECT id_kereta, harga_tiket, nama_kereta, tujuan, asal FROM public.jadwal"
		rows, err := db.Query(SQL)
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, &JadwalResponse{
				Code:   http.StatusBadRequest,
				Status: err.Error(),
			})
		}

		var result []*JadwalModel

		for rows.Next() {
			jadwal := new(JadwalModel)
			rows.Scan(&jadwal.IdKereta, &jadwal.HargaTiket, &jadwal.NamaKereta, &jadwal.Tujuan, &jadwal.Asal)

			result = append(result, jadwal)
		}

		return ctx.JSON(http.StatusOK, &JadwalResponse{
			Code:   http.StatusOK,
			Status: "success",
			Data:   result,
		})
	})

	e.POST("/tickets", func(ctx echo.Context) error {
		type TicketsRequest struct {
			IdTx     string `json:"idTx"`
			IdKereta string `json:"idKereta"`
			Nik      string `json:"nik"`
			Tujuan   string `json:"tujuan"`
			Asal     string `json:"asal"`
			Total    int    `json:"total"`
			Tanggal  string `json:"tanggal"`
		}

		type TicketsResponse struct {
			Code   int    `json:"code"`
			Status string `json:"status"`
		}

		// get body request
		req := new(TicketsRequest)
		err := ctx.Bind(req)
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, &TicketsResponse{
				Code:   http.StatusBadRequest,
				Status: err.Error(),
			})
		}

		randomNumber := rand.Intn(500)
		req.IdTx = "TX-" + strconv.Itoa(randomNumber)

		// add data to customer
		SQL := "INSERT INTO public.transaksi (id_tx, id_kereta, nik, tujuan, asal, total, tanggal) VALUES ($1, $2, $3, $4, $5, $6, $7)"
		_, err = db.Exec(SQL, req.IdTx, req.IdKereta, req.Nik, req.Tujuan, req.Asal, req.Total, req.Tanggal)
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, &TicketsResponse{
				Code: http.StatusBadRequest})
		}

		return ctx.JSON(http.StatusOK, &TicketsResponse{
			Code:   http.StatusOK,
			Status: "success",
		})
	})

	e.GET("/orders/:nik", func(ctx echo.Context) error {
		nik := ctx.Param("nik")

		type OrdersModel struct {
			IdTx     string `json:"idTx"`
			IdKereta string `json:"idKereta"`
			Nik      string `json:"nik"`
			Tujuan   string `json:"tujuan"`
			Asal     string `json:"asal"`
			Total    int    `json:"total"`
			Tanggal  string `json:"tanggal"`
		}

		type OrdersResponse struct {
			Code   int            `json:"code"`
			Status string         `json:"status"`
			Data   []*OrdersModel `json:"data"`
		}

		var result []*OrdersModel

		SQL := "SELECT id_tx, id_kereta, nik, tujuan, asal, total, tanggal FROM public.transaksi WHERE nik = $1"
		rows, err := db.Query(SQL, nik)
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, &OrdersResponse{
				Code:   http.StatusBadRequest,
				Status: err.Error(),
			})
		}

		for rows.Next() {
			orders := new(OrdersModel)
			err := rows.Scan(&orders.IdTx, &orders.IdKereta, &orders.Nik, &orders.Tujuan, &orders.Asal, &orders.Total, &orders.Tanggal)
			if err != nil {
				return ctx.JSON(http.StatusBadRequest, &OrdersResponse{
					Code:   http.StatusBadRequest,
					Status: err.Error(),
				})
			}

			result = append(result, orders)
		}

		return ctx.JSON(http.StatusOK, &OrdersResponse{
			Code:   http.StatusOK,
			Status: "success",
			Data:   result,
		})
	})

	e.POST("/jadwal", func(ctx echo.Context) error {
		type JadwalRequest struct {
			IdKereta   string `json:"idKereta"`
			HargaTiket int    `json:"hargaTiket"`
			NamaKereta string `json:"namaKereta"`
			Tujuan     string `json:"tujuan"`
			Asal       string `json:"asal"`
		}

		type JadwalResponse struct {
			Code   int            `json:"code"`
			Status string         `json:"status"`
			Data   *JadwalRequest `json:"data"`
		}

		// get body request
		req := new(JadwalRequest)
		err := ctx.Bind(req)
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, &JadwalResponse{
				Code:   http.StatusBadRequest,
				Status: err.Error(),
			})
		}

		SQL := "INSERT INTO public.jadwal (id_kereta, harga_tiket, nama_kereta, tujuan, asal) VALUES ($1, $2, $3, $4, $5)"
		_, err = db.Exec(SQL, req.IdKereta, req.HargaTiket, req.NamaKereta, req.Tujuan, req.Asal)
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, &JadwalResponse{
				Code:   http.StatusBadRequest,
				Status: err.Error(),
			})
		}

		return ctx.JSON(http.StatusOK, &JadwalResponse{
			Code:   http.StatusOK,
			Status: "success",
			Data:   req,
		})
	})

	e.PUT("/jadwal", func(ctx echo.Context) error {
		type JadwalRequest struct {
			IdKereta   string `json:"idKereta"`
			HargaTiket int    `json:"hargaTiket"`
			NamaKereta string `json:"namaKereta"`
			Tujuan     string `json:"tujuan"`
			Asal       string `json:"asal"`
		}

		type JadwalResponse struct {
			Code   int            `json:"code"`
			Status string         `json:"status"`
			Data   *JadwalRequest `json:"data"`
		}

		// get body request
		req := new(JadwalRequest)
		err := ctx.Bind(req)
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, &JadwalResponse{
				Code:   http.StatusBadRequest,
				Status: err.Error(),
			})
		}

		SQL := "UPDATE public.jadwal SET harga_tiket = $1, nama_kereta = $2, tujuan = $3, asal = $4 WHERE id_kereta = $5"
		_, err = db.Exec(SQL, req.HargaTiket, req.NamaKereta, req.Tujuan, req.Asal, req.IdKereta)
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, &JadwalResponse{
				Code:   http.StatusBadRequest,
				Status: err.Error(),
			})
		}

		return ctx.JSON(http.StatusOK, &JadwalResponse{
			Code:   http.StatusOK,
			Status: "success",
			Data:   req,
		})
	})

	e.DELETE("/jadwal/:id", func(ctx echo.Context) error {
		id := ctx.Param("id")
		SQL := "DELETE FROM public.jadwal WHERE id_kereta = ($1)"

		type Response struct {
			Code   int    `json:"code"`
			String string `json:"status"`
		}

		_, err := db.Exec(SQL, id)
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, &Response{
				Code:   http.StatusInternalServerError,
				String: err.Error(),
			})
		}

		return ctx.JSON(http.StatusOK, &Response{
			Code:   http.StatusOK,
			String: "success delete data",
		})
	})

	e.GET("user/:email", func(ctx echo.Context) error {
		type User struct {
			Nik   string `json:"nik"`
			Nama  string `json:"nama"`
			Email string `json:"email"`
			NoHp  string `json:"noHp"`
		}

		type Respose struct {
			Code   int    `json:"code"`
			Status string `json:"status"`
			Data   *User  `json:"data"`
		}

		email := ctx.Param("email")

		SQL := "SELECT nik, email, nama, no_hp FROM public.customer WHERE email = $1"
		rows, err := db.Query(SQL, email)
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, &Respose{
				Code:   http.StatusBadRequest,
				Status: err.Error(),
			})
		}

		var user User
		if rows.Next() {
			err = rows.Scan(&user.Nik, &user.Email, &user.Nama, &user.NoHp)
			if err != nil {
				return ctx.JSON(http.StatusBadRequest, &Respose{
					Code:   http.StatusBadRequest,
					Status: err.Error(),
				})
			}
		}

		return ctx.JSON(http.StatusOK, &Respose{
			Code:   http.StatusOK,
			Status: "success",
			Data:   &user,
		})
	})

	port := ":" + os.Getenv("PORT")
	e.Logger.Fatal(e.Start(port))
}
