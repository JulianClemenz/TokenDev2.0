package utils

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) { //metodo para los registros de usuarios, encriptando passwords
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool { //metodo para el ingreso a una cuenta, comparando password con hash
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	return err == nil
}
