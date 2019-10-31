package jwt

import (
	"crypto/rsa"
	"encoding"
	"math/big"
	"testing"
)

var rsaPublicKey1 *rsa.PublicKey
var rsaPublicKey2 *rsa.PublicKey
var rsaPrivateKey1 *rsa.PrivateKey
var rsaPrivateKey2 *rsa.PrivateKey

func init() {
	fromBase10 := func(base10 string) *big.Int {
		i, _ := new(big.Int).SetString(base10, 10)
		return i
	}

	rsaPublicKey1 = &rsa.PublicKey{
		N: fromBase10("14314132931241006650998084889274020608918049032671858325988396851334124245188214251956198731333464217832226406088020736932173064754214329009979944037640912127943488972644697423190955557435910767690712778463524983667852819010259499695177313115447116110358524558307947613422897787329221478860907963827160223559690523660574329011927531289655711860504630573766609239332569210831325633840174683944553667352219670930408593321661375473885147973879086994006440025257225431977751512374815915392249179976902953721486040787792801849818254465486633791826766873076617116727073077821584676715609985777563958286637185868165868520557"),
		E: 3,
	}
	rsaPrivateKey1 = &rsa.PrivateKey{
		PublicKey: *rsaPublicKey1,
		D:         fromBase10("9542755287494004433998723259516013739278699355114572217325597900889416163458809501304132487555642811888150937392013824621448709836142886006653296025093941418628992648429798282127303704957273845127141852309016655778568546006839666463451542076964744073572349705538631742281931858219480985907271975884773482372966847639853897890615456605598071088189838676728836833012254065983259638538107719766738032720239892094196108713378822882383694456030043492571063441943847195939549773271694647657549658603365629458610273821292232646334717612674519997533901052790334279661754176490593041941863932308687197618671528035670452762731"),
		Primes: []*big.Int{
			fromBase10("130903255182996722426771613606077755295583329135067340152947172868415809027537376306193179624298874215608270802054347609836776473930072411958753044562214537013874103802006369634761074377213995983876788718033850153719421695468704276694983032644416930879093914927146648402139231293035971427838068945045019075433"),
			fromBase10("109348945610485453577574767652527472924289229538286649661240938988020367005475727988253438647560958573506159449538793540472829815903949343191091817779240101054552748665267574271163617694640513549693841337820602726596756351006149518830932261246698766355347898158548465400674856021497190430791824869615170301029"),
		},
	}

	rsaPublicKey2 = &rsa.PublicKey{
		N: fromBase10("887256638780856047914581579814082241228904817182093261763058947771445376334361454028526564964284108022499562499436107360797167855238681264022398768822093355881651855771536641100347179464316604983028325045398953751392004347747368839824636917782721516826693522992183854650791268736264416865889799352514216298521461994511678593557464913439951619056304159812403323561755303818800310137159146108075344562143531014989113550013919515074346824023514684979514198941230125010212380546452298362548591023662125669118549327451011201031961777374564883854203295916398755760815894642777422649649856180524299517641811337173222323171868"),
		E: 3,
	}
	rsaPrivateKey2 = &rsa.PrivateKey{
		PublicKey: *rsaPublicKey2,
		D:         fromBase10("9542755287494004433998723259516013739278699355114572217325597900889416163458809501304132487555642811888150937392013824621448709836142886006653296025093941418628992648429798282127303704957273845127141852309016655778568546006839666463451542076964744073572349705538631742281931858219480985907271975884773482372966847639853897890615456605598071088189838676728836833012254065983259638538107719766738032720239892094196108713378822882383694456030043492571063441943847195939549773271694647657549658603365629458610273821292232646334717612674519997533901052790334279661754176490593041941863932308687197618671528035670452762731"),
		Primes: []*big.Int{
			fromBase10("118824431543227186571027026092363352199324615266444280871852738808378452065820764461459770308550270192458287541328152308596233985238045601785765785685887955407926672604456991769558426224026688264520894241180697554536183714315783324715006626168049073199347594517685254572389805736483316580732446129950642427691"),
			fromBase10("167151560084187340876600839609071108780428930679164316900248316388276108158369650170260479635956969047912292329416883197284109866041024250716279228267907058587184812901048552331276732325276799660817363854674774994782375599727993270112962045907225952829559312944081207881075710161619024002861307509427534745757"),
		},
	}
}

func TestRS256_WithValidSignature(t *testing.T) {
	f := func(signer Signer, claims encoding.BinaryMarshaler) {
		t.Helper()

		tokenBuilder := NewTokenBuilder(signer)
		token, _ := tokenBuilder.Build(claims)

		err := signer.Verify(token.Payload(), token.Signature())
		if err != nil {
			t.Errorf("want no err, got: `%v`", err)
		}
	}

	f(
		NewRS256(rsaPublicKey1, rsaPrivateKey1),
		&StandardClaims{},
	)
	f(
		NewRS384(rsaPublicKey1, rsaPrivateKey1),
		&StandardClaims{},
	)
	f(
		NewRS512(rsaPublicKey1, rsaPrivateKey1),
		&StandardClaims{},
	)

	f(
		NewRS256(rsaPublicKey1, rsaPrivateKey1),
		&customClaims{
			TestField: "foo",
		},
	)
	f(
		NewRS384(rsaPublicKey1, rsaPrivateKey1),
		&customClaims{
			TestField: "bar",
		},
	)
	f(
		NewRS512(rsaPublicKey1, rsaPrivateKey1),
		&customClaims{
			TestField: "baz",
		},
	)
}

func TestRS384_WithInvalidSignature(t *testing.T) {
	f := func(signer, verifier Signer, claims encoding.BinaryMarshaler) {
		t.Helper()

		tokenBuilder := NewTokenBuilder(signer)
		token, _ := tokenBuilder.Build(claims)

		err := verifier.Verify(token.Payload(), token.Signature())
		if err == nil {
			t.Errorf("want %v, got nil", ErrInvalidSignature)
		}
	}
	f(
		NewRS256(rsaPublicKey1, rsaPrivateKey1),
		NewRS256(rsaPublicKey2, rsaPrivateKey2),
		&StandardClaims{},
	)
	f(
		NewRS384(rsaPublicKey1, rsaPrivateKey1),
		NewRS384(rsaPublicKey2, rsaPrivateKey2),
		&StandardClaims{},
	)
	f(
		NewRS512(rsaPublicKey1, rsaPrivateKey1),
		NewRS512(rsaPublicKey2, rsaPrivateKey2),
		&StandardClaims{},
	)

	f(
		NewRS256(rsaPublicKey1, rsaPrivateKey1),
		NewRS256(rsaPublicKey2, rsaPrivateKey2),
		&customClaims{
			TestField: "foo",
		},
	)
	f(
		NewRS384(rsaPublicKey1, rsaPrivateKey1),
		NewRS384(rsaPublicKey2, rsaPrivateKey2),
		&customClaims{
			TestField: "bar",
		},
	)
	f(
		NewRS512(rsaPublicKey1, rsaPrivateKey1),
		NewRS512(rsaPublicKey2, rsaPrivateKey2),
		&customClaims{
			TestField: "baz",
		},
	)
}
