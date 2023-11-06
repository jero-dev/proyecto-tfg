package services_test

import (
	"testing"
	services "vidya-sale/api/service"
)

func TestMessageProcessorService_ParseMessage(t *testing.T) {
	var messageProcessor services.MessageProcessorService

	type testCase struct {
		name         string
		message      string
		expectedData struct {
			gameName string
			platform string
			link     string
			price    float64
		}
	}

	testCases := []testCase{
		{
			name: "Example 1",
			message: `⬇️ Red Dead Redemption #Switch
            (https://i.ibb.co/4RDsKSH/ofertasjuegoses.png)BAJONAZO en reserva de 50 a 44,99€, mínimo alcanzado (PVP 50€) (Distribuye Nintendo, así que no va a bajar mucho en tiempo)
    
            🟣 https://amzn.to/446ayBA
        	`,
			expectedData: struct {
				gameName string
				platform string
				link     string
				price    float64
			}{
				gameName: "Red Dead Redemption",
				platform: "Switch",
				link:     "https://amzn.to/446ayBA",
				price:    44.99,
			},
		},
		{
			name: "Example 2",
			message: `⬇️ Detective Pikachu: El Regreso #Switch
			(https://telegra.ph/file/c2562f5797578dd3db025.png)BAJONAZO SUPREMO FLASH a solo 34,29€ ¡Su nuevo precio mínimo sin cupones! (PVP 50€)
		   ✂️ Con el nuevo cupón de primera compra (https://t.me/OfertasJuegosNintendo/20221) se queda a solo 24€, nuevo mínimo total
		   
		   🌸 https://ojueg.es/hGXg4
		   
		   📝 Edición española con envíos rápidos de parte de Gamers4Life
		   
		   ⚪️ Si no te aclaras con el tema del cupón, averigua como sacarle partido en nuestro hilo de dudas Miravia del @GrupoNintendoOJ 
		   🍨 https://t.me/GrupoNintendoOJ/481642
			`,
			expectedData: struct {
				gameName string
				platform string
				link     string
				price    float64
			}{
				gameName: "Detective Pikachu: El Regreso",
				platform: "Switch",
				link:     "https://ojueg.es/hGXg4",
				price:    34.29,
			},
		},
		{
			name: "Example 3",
			message: `📆⬇️ Star Ocean The Second Story R #Switch
			(https://telegra.ph/file/2c7ec774526506600df7a.png)BAJONAZO FLASH a solo 47,51€, precio mínimo alcanzado incluso sin cupón (PVP 60€)
		   ✂️ Con cupón de primera compra se te queda a solo 33,25€
		   
		   🌸 https://ojueg.es/vgfPa
		   
		   
		   📝 Edición española con envío de lanzamiento de parte de Gamers4Life
		   
		   
		   ⚪️ Si no te aclaras con el tema del cupón, averigua como sacarle partido en nuestro hilo de dudas Miravia del @GrupoNintendoOJ 
		   🍨 https://t.me/GrupoNintendoOJ/481642
			`,
			expectedData: struct {
				gameName string
				platform string
				link     string
				price    float64
			}{
				gameName: "Star Ocean The Second Story R",
				platform: "Switch",
				link:     "https://ojueg.es/vgfPa",
				price:    47.51,
			},
		},
		{
			name: "Example 4",
			message: `📆⬇️ Another Code Recollection #Switch (https://telegra.ph/file/66470e955a21aa6f88d6e.png)
				BAJONAZO SUPREMO en reserva de 55 a 46,99€ en Amazon, su mínimo alcanzado (PVP 60€)
				🟣 https://amzn.to/3PWgadS
		
				❤ Recuerda que en Amazon no se paga nada por reservar, si bajara incluso más te guardan su precio mínimo y se podría cancelar en un click si lo deseas.
		
				✅️ El próximo en acomodar un precio mega-top será el Mario VS Donkey Kong (PVP 50€). Si lo reservas ya, todo apunta que se quedará a 39,90€ y te lo igualarán automáticamente
				🟣 https://amzn.to/3RgP7LA
			`,
			expectedData: struct {
				gameName string
				platform string
				link     string
				price    float64
			}{
				gameName: "Another Code Recollection",
				platform: "Switch",
				link:     "https://amzn.to/3PWgadS",
				price:    46.99,
			},
		},
		{
			name: "Example 5",
			message: `⬇️ Super Mario Bros Wonder #Switch (https://telegra.ph/file/9d92a01d78575158a659f.png)
				(https://telegra.ph/file/789f6336bb488f3c20cef.png)NUEVA OPORTUNIDAD FLASH a solo 46,80€
				✂️ Se queda a solo 32,76€ usando cupón de primera compra
				🌸 https://ojueg.es/enIls
		
				📝 Edición española,  vendido por TheShopGamer con envíos rápidos
		
				✅️ Miravia nos ofrece diariamente mejores ofertas que el día sin IVA, echa un vistazo
				🌸  https://ojueg.es/MiraviaFlash
		
				⚪️Si tienes alguna duda, problemas creando la cuenta, no te sale el cupón, etc, vente al hilo de dudas en nuestro @GrupoNintendoOJ para echar un cable
				🍨 https://t.me/GrupoNintendoOJ/481642
			`,
			expectedData: struct {
				gameName string
				platform string
				link     string
				price    float64
			}{
				gameName: "Super Mario Bros Wonder",
				platform: "Switch",
				link:     "https://ojueg.es/enIls",
				price:    46.80,
			},
		},
		{
			name: "Example 6",
			message: `⬇️ Mario Kart 8 Booster Pack #Switch (https://telegra.ph/file/1302008dd18752c139cd7.png)
				BAJONAZO SUPREMO FLASH a solo 26,96€ (PVP 40€)
				✂️ Usando el cupón de primera compra se queda a solo 18,87€
				🌸 https://ojueg.es/BlbhN
		
				📝 Edición española vendida por PreciosYa con envíos rápidos. Teniendo en cuenta que la expansión cuesta 25€ directamente en digital (Bueno, 19,20€ en IG (https://ojueg.es/UxqXv)), te llevas el set de merchandising de regalo y el DLC a precio top.
		
				⚪️ Si no te aclaras con el tema del cupón, averigua como sacarle partido en nuestro hilo de dudas Miravia del @GrupoNintendoOJ 
				🍨 https://t.me/GrupoNintendoOJ/481642
			`,
			expectedData: struct {
				gameName string
				platform string
				link     string
				price    float64
			}{
				gameName: "Mario Kart 8 Booster Pack",
				platform: "Switch",
				link:     "https://ojueg.es/BlbhN",
				price:    26.96,
			},
		},
		{
			name: "Example 7",
			message: `📆⬇️ Hogwarts Legacy #Switch
				 (https://telegra.ph/file/68bbe427dd8d7f14fdec2.png)BAJONAZO SUPREMO FLASH a solo 40,30€ (PVP 60€)
				✂️ Con cupón de primera compra se queda a solo 28,21€
				🌸 https://ojueg.es/GEIRp
		
				📝 Edición española vendida por PreciosYa para recibirla de lanzamiento el 14 de noviembre
		
				✅️ Miravia nos ofrece diariamente mejores ofertas que el día sin IVA, echa un vistazo
				🌸  https://ojueg.es/MiraviaFlash
		
				⚪️ Si no te aclaras con el tema del cupón, averigua como sacarle partido en nuestro hilo de dudas Miravia del @GrupoNintendoOJ 
				🍨 https://t.me/GrupoNintendoOJ/481642
			`,
			expectedData: struct {
				gameName string
				platform string
				link     string
				price    float64
			}{
				gameName: "Hogwarts Legacy",
				platform: "Switch",
				link:     "https://ojueg.es/GEIRp",
				price:    40.30,
			},
		},
		{
			name: "Example 8",
			message: `⬇️ Super Mario Wonder #Switch 
			 (https://telegra.ph/file/d533f9f55dd5fe20c8c41.png)BAJONAZO SUPREMO FLASH a solo 44,80€, nuevo precio mínimo  ¡Hasta las 00h o fin de existencias! (PVP 60€)
			✂️Usando cupón de primera compra (https://t.me/OfertasJuegosNintendo/20221) queda a solo 31,36€, nuevo mínimo total
			🌸 https://ojueg.es/3RTs1
		
			📝Edición española vendida por PreciosYa con envíos rápidos
		
			✅️ Miravia nos ofrece diariamente mejores ofertas que el día sin IVA, echa un vistazo
			🌸  https://ojueg.es/MiraviaFlash
		
			⚪️ Si no te aclaras con el tema del cupón, averigua como sacarle partido en nuestro hilo de dudas Miravia del @GrupoNintendoOJ 
			🍨 https://t.me/GrupoNintendoOJ/481642
			`,
			expectedData: struct {
				gameName string
				platform string
				link     string
				price    float64
			}{
				gameName: "Super Mario Wonder",
				platform: "Switch",
				link:     "https://ojueg.es/3RTs1",
				price:    44.80,
			},
		},
		{
			name: "Example 9",
			message: `📆 ⬇️ Super Mario RPG #Switch
			 (https://telegra.ph/file/529625f18b52c1548d6b5.png)BAJONAZO SUPREMO FLASH en reserva a solo 44,98€, nuevo precio mínimo incluso sin cupón (PVP 60€)
			✂️ Usando el cupón de primera compra (https://t.me/ofertasjuegoses/44974) a solo 31,48€, nuevo precio mínimo total en reserva
			🌸 https://ojueg.es/D1DrA
		
			📝 Edición española y de lanzamiento vendido por PreciosYa
		
			⚪️Si tienes alguna duda, problemas creando la cuenta, no te sale el cupón, etc, vente al hilo de dudas en nuestro @GrupoNintendoOJ para echar un cable
			🍨 https://t.me/GrupoNintendoOJ/481642
			`,
			expectedData: struct {
				gameName string
				platform string
				link     string
				price    float64
			}{
				gameName: "Super Mario RPG",
				platform: "Switch",
				link:     "https://ojueg.es/D1DrA",
				price:    44.98,
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			gameName, platform, link, price := messageProcessor.ParseMessage(testCase.message)

			if gameName != testCase.expectedData.gameName {
				t.Errorf("Expected game name to be %v but got %v", testCase.expectedData.gameName, gameName)
			}
			if platform != testCase.expectedData.platform {
				t.Errorf("Expected platform to be %v but got %v", testCase.expectedData.platform, platform)
			}
			if link != testCase.expectedData.link {
				t.Errorf("Expected link to be %v but got %v", testCase.expectedData.link, link)
			}
			if price != testCase.expectedData.price {
				t.Errorf("Expected price to be %v but got %v", testCase.expectedData.price, price)
			}
		})
	}
}
