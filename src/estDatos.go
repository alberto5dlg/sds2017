package main

type logueado struct {
	User     string
	Password string
}

type cuenta struct {
	Boss     string
	Servicio string
	User     string
	Password string
}

type registrarse struct {
	User     string
	Password string
	Email    string
}

type datos struct {
	User string
	Pass string
}

type usuario struct {
	Email    string
	Password string
	Info     map[string]datos
	Tarjetas map[string]tarjeta
	Notas    map[string]notas
}

type resp struct {
	Ok  bool
	Msg string
}

type respJSON struct {
	Ok   bool
	Info map[string]datos
}

type nTarjeta struct {
	Username string
	Entidad  string
	NTarjeta string
	Fecha    string
	CodSeg   string
}

type tarjeta struct {
	Entidad  string
	NTarjeta string
	Fecha    string
	CodSeg   string
}

type nNotas struct {
	Username string
	Titulo   string
	Cuerpo   string
}

type notas struct {
	Titulo string
	Cuerpo string
}
