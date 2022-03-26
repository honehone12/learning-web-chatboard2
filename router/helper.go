package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"learning-web-chatboard2/common"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	runeSource         = "aA1bB2cC3dD4eE5fFgGhHiIjJkKlLm0MnNoOpPqQrRsStTuUvV6wW7xX8yY9zZ"
	macSalt            = "uPUqL7dZ"
	pwSalt             = "LV2vP8vq"
	sessionCookieLabel = "short-time"
	visitCookieLabel   = "long-time"
)
const (
	aes256KeySize uint          = 32
	macKeySize    uint          = 32
	stateSize     uint          = 32
	sessionExp    time.Duration = time.Hour * 8
	stateExp      time.Duration = time.Minute * 20
	visitExp      time.Duration = time.Hour * 24 * 365
)

var helper struct {
	block  cipher.Block
	macKey []byte
}

// every time server is restarted, cookie become no longer valid
func startHelper() (err error) {
	bKyeStr, err := generateString(aes256KeySize)
	if err != nil {
		return
	}
	bKey := []byte(bKyeStr)
	macKeyStr, err := generateString(macKeySize)
	if err != nil {
		return
	}
	helper.macKey = []byte(macKeyStr)
	helper.block, err = aes.NewCipher(bKey)
	return
}

func buildHTTP_URL(domain string, path string) (url string) {
	url = fmt.Sprintf(
		"%s%s%s",
		httpPrefix,
		domain,
		path,
	)
	return
}

func makeHash(plainText string) (hashed string) {
	asBytes := sha256.Sum256([]byte(plainText))
	hashed = fmt.Sprintf("%x", asBytes)
	return
}

func processPassword(pw string) string {
	// see these pkgs
	// https://pkg.go.dev/golang.org/x/crypto/bcrypt
	// https://pkg.go.dev/golang.org/x/crypto/scrypt
	return makeHash(fmt.Sprint(pwSalt, pw))
}

func generateString(length uint) (str string, err error) {
	var i uint
	maxEx := int64(len(runeSource))
	runePool := []rune(runeSource)
	for i = 0; i < length; i++ {
		bigN, err := rand.Int(rand.Reader, big.NewInt(maxEx))
		if err != nil {
			break
		}
		n := bigN.Uint64()
		str = fmt.Sprint(str, string(runePool[n]))
	}
	return
}

func encrypt(plainText string) (cipherText []byte, err error) {
	cipherText = make([]byte, aes.BlockSize+len(plainText))
	iv := cipherText[:aes.BlockSize]
	n, err := io.ReadFull(rand.Reader, iv)
	if err != nil {
		err = fmt.Errorf("%s: returned %d", err.Error(), n)
		return
	}

	encryptStream := cipher.NewCTR(helper.block, iv)
	encryptStream.XORKeyStream(cipherText[aes.BlockSize:], []byte(plainText))
	return
}

func decrypt(cipherText []byte) (plainText string, err error) {
	decryptText := make([]byte, len(cipherText[aes.BlockSize:]))
	decryptStream := cipher.NewCTR(helper.block, cipherText[:aes.BlockSize])
	decryptStream.XORKeyStream(decryptText, cipherText[aes.BlockSize:])
	plainText = string(decryptText)
	return
}

func makeMAC(value []byte) []byte {
	hash := hmac.New(sha256.New, helper.macKey)
	hash.Write(value)
	return hash.Sum([]byte(macSalt))
}

func verifyMAC(mac []byte, value []byte) bool {
	hash := hmac.New(sha256.New, helper.macKey)
	hash.Write(value)
	hashedVal := hash.Sum([]byte(macSalt))
	return hmac.Equal(mac, hashedVal)
}

func encode(value []byte) string {
	return base64.URLEncoding.EncodeToString(value)
}

func decode(encoded string) ([]byte, error) {
	return base64.URLEncoding.DecodeString(encoded)
}

func checkLoggedIn(ctx *gin.Context) (err error) {
	uuid, err := pickupCookie(ctx, sessionCookieLabel)
	if err != nil {
		return
	}
	sess := common.Session{
		UuId: uuid,
	}
	req, err := common.MakeRequestFromSession(
		&sess,
		http.MethodPost,
		buildHTTP_URL(config.AddressUsers, "/check-session"),
	)
	if err != nil {
		return
	}
	res, err := httpClient.Do(req)
	if err != nil {
		return
	} else if res.StatusCode != http.StatusOK {
		err = errors.New(res.Status)
		return
	}
	sessPtr, err := common.MakeSessionFromResponse(res)
	if err != nil {
		return
	}

	ctx.Set(sessionPtrLabel, sessPtr)
	return
}

func visitCheck(ctx *gin.Context) (err error) {
	var vis *common.Visit
	_, err = pickupCookie(ctx, visitCookieLabel)
	if err == nil {
		vis, err = requestVisitPtr(ctx)
	}
	if err != nil {
		if gin.IsDebugging() {
			common.LogWarning(logger).
				Printf("creating new visit because [%s]\n", err.Error())
		}
		vis, err = requestVisitCreate()
		if err != nil {
			return
		}
		storeVisitCookie(ctx, vis.UuId)
	}

	ctx.Set(visitPtrLabel, vis)
	err = nil
	return
}

func storeSessionCookie(ctx *gin.Context, value string) (err error) {
	err = storeCookie(
		ctx,
		value,
		sessionCookieLabel,
		sessionExp,
		0,
	)
	return
}

func storeVisitCookie(ctx *gin.Context, value string) (err error) {
	err = storeCookie(
		ctx,
		value,
		visitCookieLabel,
		visitExp,
		60*60*24*365,
	)
	return
}

func storeCookie(
	ctx *gin.Context,
	value string,
	cookieName string,
	sessionDuration time.Duration,
	cookieDuration int,
) (err error) {
	//add exp
	value = fmt.Sprintf(
		"%s|%d",
		value,
		time.Now().Add(sessionDuration).Unix(),
	)
	encrypted, err := encrypt(value)
	if err != nil {
		return
	}
	// add mac value first
	bytesVal := makeMAC(encrypted)
	// separated '|'
	bytesVal = append(bytesVal, []byte("|")...)
	// add encrypted value
	bytesVal = append(bytesVal, encrypted...)

	valToStore := encode(bytesVal)

	if gin.IsDebugging() {
		common.LogInfo(logger).
			Printf("stored cookie [%s] %s\n", cookieName, valToStore)
	}
	ctx.SetSameSite(http.SameSiteStrictMode)
	ctx.SetCookie(
		cookieName,
		valToStore,
		cookieDuration,
		"/",
		config.AddressRouter,
		true,
		true,
	)
	return
}

func pickupCookie(ctx *gin.Context, name string) (value string, err error) {
	rawStored, err := ctx.Cookie(name)
	if err != nil {
		return
	}
	bytesVal, err := decode(rawStored)
	if err != nil {
		return
	}
	splited := bytes.SplitN(bytesVal, []byte("|"), 2)
	mac := splited[0]
	encrypted := splited[1]
	if !verifyMAC(mac, encrypted) {
		err = fmt.Errorf("invalid cookie %s", rawStored)
		return
	}
	decrypted, err := decrypt(encrypted)
	if err != nil {
		return
	}
	value, unixTimeStr, ok := strings.Cut(decrypted, "|")
	if !ok {
		err = errors.New("separator not found")
		return
	}
	unixTime, err := strconv.ParseInt(unixTimeStr, 10, 64)
	if err != nil {
		return
	}

	if unixTime < time.Now().Unix() {
		err = errors.New("session expired")
	}
	return
}

func generateVisitState(ctx *gin.Context) (stateAndMACEncoded string, err error) {
	vis, err := requestVisitPtr(ctx)
	if err != nil {
		return
	}

	vis.State, stateAndMACEncoded, err = generateState()
	if err != nil {
		return
	}
	err = requestVisitUpdate(vis)
	if err != nil {
		return
	}
	ctx.Set(visitPtrLabel, vis)
	return
}

func generateSessionState(ctx *gin.Context) (stateAndMACEncoded string, err error) {
	sess, err := getSessionPtrFromCTX(ctx)
	if err != nil {
		return
	}

	sess.State, stateAndMACEncoded, err = generateState()
	if err != nil {
		return
	}
	err = requestSessionUpdate(sess)
	if err != nil {
		return
	}
	ctx.Set(sessionPtrLabel, sess)
	return
}

func generateState() (stateRaw, stateAndMACEncoded string, err error) {
	state, err := generateString(stateSize)
	if err != nil {
		return
	}
	state = fmt.Sprintf(
		"%s|%d",
		state,
		time.Now().Add(stateExp).Unix(),
	)
	stateRaw = state

	// same proc with cookie
	stateAsBytes := []byte(state)
	bytesVal := makeMAC(stateAsBytes)
	bytesVal = append(bytesVal, []byte("|")...)
	bytesVal = append(bytesVal, stateAsBytes...)
	stateAndMACEncoded = encode(bytesVal)
	return
}

func requestSessionUpdate(sess *common.Session) (err error) {
	req, err := common.MakeRequestFromSession(
		sess,
		http.MethodPost,
		buildHTTP_URL(config.AddressUsers, "/update-session"),
	)
	if err != nil {
		return
	}
	res, err := httpClient.Do(req)
	if err != nil {
		return
	} else if res.StatusCode != http.StatusOK {
		err = errors.New(res.Status)
		return
	}
	sess, err = common.MakeSessionFromResponse(res)
	return
}

func requestVisitCreate() (vis *common.Visit, err error) {
	req, err := http.NewRequest(
		http.MethodGet,
		buildHTTP_URL(config.AddressUsers, "/create-visit"),
		nil,
	)
	if err != nil {
		return
	}
	res, err := httpClient.Do(req)
	if err != nil {
		return
	} else if res.StatusCode != http.StatusOK {
		err = errors.New(res.Status)
		return
	}
	vis, err = common.MakeVisitFromResponse(res)
	return
}

func requestVisitPtr(ctx *gin.Context) (ptr *common.Visit, err error) {
	uuid, err := pickupCookie(ctx, visitCookieLabel)
	if err != nil {
		return
	}
	vis := common.Visit{UuId: uuid}
	req, err := common.MakeRequestFromVisit(
		&vis,
		http.MethodPost,
		buildHTTP_URL(config.AddressUsers, "/check-visit"),
	)
	if err != nil {
		return
	}
	res, err := httpClient.Do(req)
	if err != nil {
		return
	} else if res.StatusCode != http.StatusOK {
		err = errors.New(res.Status)
		return
	}
	ptr, err = common.MakeVisitFromResponse(res)
	return
}

func requestVisitUpdate(vis *common.Visit) (err error) {
	req, err := common.MakeRequestFromVisit(
		vis,
		http.MethodPost,
		buildHTTP_URL(config.AddressUsers, "/update-visit"),
	)
	if err != nil {
		return
	}
	res, err := httpClient.Do(req)
	if err != nil {
		return
	} else if res.StatusCode != http.StatusOK {
		err = errors.New(res.Status)
	}
	return
}

func checkState(exposedVal, privateVal string) (err error) {
	if strings.Compare(exposedVal, "") == 0 {
		err = errors.New("exposed value is empty")
		return
	}
	if strings.Compare(privateVal, "") == 0 {
		err = errors.New("private value is empty")
		return
	}

	bytesVal, err := decode(exposedVal)
	if err != nil {
		return
	}
	splited := bytes.SplitN(bytesVal, []byte("|"), 2)
	macStored := splited[0]
	stateStored := string(splited[1])

	if !verifyMAC(macStored, []byte(privateVal)) {
		err = errors.New("invalid mac")
		return
	}
	if strings.Compare(stateStored, privateVal) != 0 {
		err = errors.New("invalid state")
		return
	}
	_, unixTimeStr, ok := strings.Cut(stateStored, "|")
	if !ok {
		err = errors.New("separator not found")
		return
	}
	unixTime, err := strconv.ParseInt(unixTimeStr, 10, 64)
	if err != nil {
		return
	}
	if unixTime < time.Now().Unix() {
		err = errors.New("state expired")
	}
	return
}

func visitStateCheckProcess(ctx *gin.Context) (vis *common.Visit, err error) {
	vis, err = getVisitPtrFromCTX(ctx)
	if err != nil {
		return
	}

	// check state
	state := ctx.PostForm("state")
	err = checkState(state, vis.State)
	if err != nil {
		return
	}

	// state is consumed, delete it
	vis.State = ""
	err = requestVisitUpdate(vis)
	return
}

func sessionStateCheckProcess(ctx *gin.Context) (sess *common.Session, err error) {
	sess, err = getSessionPtrFromCTX(ctx)
	if err != nil {
		return
	}

	// check state
	state := ctx.PostForm("state")
	err = checkState(state, sess.State)
	if err != nil {
		return
	}

	// state is consumed, delete it
	sess.State = ""
	err = requestSessionUpdate(sess)
	return
}
