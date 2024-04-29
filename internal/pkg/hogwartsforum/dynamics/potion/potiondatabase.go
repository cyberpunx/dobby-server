package potion

import "html/template"

const (
	AgudizadoraDeIngenio            = "Agudizadora de Ingenio"
	AlientoDeFuego                  = "Aliento de Fuego"
	Amortentia                      = "Amortentia"
	AntidotoParaVenenosComunes      = "Antidoto para Venenos Comunes"
	AntidotoDeGlumbumble            = "Antídoto de Glumbumble"
	AntidotoParaVenenosToxicos      = "Antídoto para Venenos Toxicos"
	BrebajeParlanchin               = "Brebaje Parlanchin"
	Doxycida                        = "Doxycida"
	ElixirParaInducirEuforia        = "Elixir para Inducir Euforia"
	Embellecedora                   = "Embellecedora"
	EsenciaDeDictamo                = "Esencia de Díctamo"
	EsenciaDeMurtlap                = "Esencia de Murtlap"
	FelixFelicis                    = "Felix Felicis"
	FertilizanteDeEstiercolDeDragon = "Fertilizante de Estiercol de Dragón"
	FiltroDeMuertosEnVida           = "Filtro de Muertos en Vida"
	FiltroDePaz                     = "Filtro de Paz"
	GraspOfDeath                    = "Grasp of Death"
	Multijugos                      = "Multijugos"
	PocionCalmante                  = "Poción Calmante"
	PocionCrecepelo                 = "Poción Crecepelo"
	PocionCuradoraDeFurunculos      = "Poción Curadora de Furúnculos"
	PocionDeDespertares             = "Poción de Despertares"
	PocionDeErumpent                = "Poción de Erumpent"
	PocionDeLaMemoria               = "Poción de la Memoria"
	PocionDeLaRisa                  = "Poción de la Risa"
	PocionDelOlvido                 = "Poción del Olvido"
	PocionHerbicida                 = "Poción Herbicida"
	PocionHerbovitalizante          = "Poción Herbovitalizante"
	PocionOculus                    = "Poción Oculus"
	PocionParaElDolorDeEstomago     = "Poción para el dolor de estómago"
	PocionPimentonica               = "Poción Pimentónica"
	PocionProtectoraContraLasLlamas = "Poción Protectora Contra las Llamas"
	PocionVigorizante               = "Poción Vigorizante"
	PocionVolubilis                 = "Poción Volubilis"
	PocimaParaDormir                = "Pócima para Dormir"
	SolucionAgrandadora             = "Solución Agrandadora"
	SolucionEncogedora              = "Solución Encogedora"
	Veritaserum                     = "Veritaserum"
	ZumoDeMandragora                = "Zumo de Mandrágora"
)

var PotionIcons = map[string]template.HTML{
	AgudizadoraDeIngenio:            `<img src="https://i.imgur.com/GanNYpr.png" class="tooltip" title="Agudizadora del Ingenio"/>`,
	Amortentia:                      `<img src="https://i.imgur.com/yuVxqDZ.png" class="tooltip" title="Amortentia"/>`,
	AntidotoDeGlumbumble:            `<img src="https://i.imgur.com/9RBZCeu.png" class="tooltip" title="Antídoto de Glumbumble"/>`,
	AntidotoParaVenenosComunes:      `<img src="https://i.imgur.com/0RRcrav.png" class="tooltip" title="Antídoto para Venenos Comunes"/>`,
	AntidotoParaVenenosToxicos:      `<img src="https://i.imgur.com/m2sCKZf.png" class="tooltip" title="Antídoto para Venenos Tóxicos"/>`,
	BrebajeParlanchin:               `<img src="https://i.imgur.com/YHbpTZb.png" class="tooltip" title="Brebaje Parlanchín"/>`,
	PocionCrecepelo:                 `<img src="https://i.imgur.com/lcAX0qO.png" class="tooltip" title="Crecepelo"/>`,
	PocionCuradoraDeFurunculos:      `<img src="https://i.imgur.com/lGHayif.png" class="tooltip" title="Curadora de Furúnculos"/>`,
	Doxycida:                        `<img src="https://i.imgur.com/ACnnCYS.png" class="tooltip" title="Doxycida"/>`,
	ElixirParaInducirEuforia:        `<img src="https://i.imgur.com/qigGa6E.png" class="tooltip" title="Elixir para inducir Euforia"/>`,
	Embellecedora:                   `<img src="https://i.imgur.com/CQSs7AL.png" class="tooltip" title="Embellecedora"/>`,
	EsenciaDeDictamo:                `<img src="https://i.imgur.com/DaWtYLH.png" class="tooltip" title="Esencia de Díctamo"/>`,
	EsenciaDeMurtlap:                `<img src="https://i.imgur.com/TPPuHVj.png" class="tooltip" title="Esencia de Murtlap"/>`,
	FelixFelicis:                    `<img src="https://i.imgur.com/o54Rv3c.png" class="tooltip" title="Felix Felicis"/>`,
	FertilizanteDeEstiercolDeDragon: `<img src="https://i.imgur.com/z7biKGy.png" class="tooltip" title="Fertilizante de Estiércol de Dragón"/>`,
	FiltroDeMuertosEnVida:           `<img src="https://i.imgur.com/HSFCLlF.png" class="tooltip" title="Filtro de Muertos en Vida"/>`,
	FiltroDePaz:                     `<img src="https://i.imgur.com/Tm4fAWN.png" class="tooltip" title="Filtro de Paz"/>`,
	GraspOfDeath:                    `<img src="https://i.imgur.com/TA2nuNv.png" class="tooltip" title="Grasp of Death"/>`,
	Multijugos:                      `<img src="https://i.imgur.com/1kEPEYJ.png" class="tooltip" title="Multijugos"/>`,
	PocimaParaDormir:                `<img src="https://i.imgur.com/Z7R6IF4.png" class="tooltip" title="Pócima para Dormir"/>`,
	AlientoDeFuego:                  `<img src="https://i.imgur.com/JiZi0Wh.png" class="tooltip" title="Poción Aliento de Fuego"/>`,
	PocionCalmante:                  `<img src="https://i.imgur.com/HTMjhyQ.png" class="tooltip" title="Poción Calmante"/>`,
	PocionDeDespertares:             `<img src="https://i.imgur.com/lfgDzmp.png" class="tooltip" title="Poción de Despertares"/>`,
	PocionDeErumpent:                `<img src="https://i.imgur.com/TzFgk4h.png" clas="tooltip" title="Poción de Erumpent"/>`,
	PocionDeLaMemoria:               `<img src="https://i.imgur.com/ldkNgXk.png" class="tooltip" title="Poción de la Memoria"/>`,
	PocionDeLaRisa:                  `<img src="https://i.imgur.com/rnrljbI.png" class="tooltip" title="Poción de la Risa"/>`,
	PocionDelOlvido:                 `<img src="https://i.imgur.com/qsTMBQd.png" class="tooltip" title="Poción del Olvido"/>`,
	PocionHerbicida:                 `<img src="https://i.imgur.com/JE6nCfo.png" class="tooltip" title="Poción Herbicida"/>`,
	PocionHerbovitalizante:          `<img src="https://i.imgur.com/FiunEWk.png" class="tooltip" title="Poción Herbovitalizante"/>`,
	PocionOculus:                    `<img src="https://i.imgur.com/0IXlBh8.png" class="tooltip" title="Poción Oculus"/>`,
	PocionParaElDolorDeEstomago:     `<img src="https://i.imgur.com/k3cFa0J.png" class="tooltip" title="Poción para el Dolor de Estómago"/>`,
	PocionPimentonica:               `<img src="https://i.imgur.com/U5zJMuB.png" class="tooltip" title="Poción Pimentónica"/>`,
	PocionProtectoraContraLasLlamas: `<img src="https://i.imgur.com/dx4Dh3y.png" class="tooltip" title="Poción Protectora contra las Llamas"/>`,
	PocionVigorizante:               `<img src="https://i.imgur.com/V5mkB3V.jpg" class="tooltip" title="Poción Vigorizante"/>`,
	PocionVolubilis:                 `<img src="https://i.imgur.com/U3SIfC6.png" class="tooltip" title="Poción Volubilis"/>`,
	SolucionAgrandadora:             `<img src="https://i.imgur.com/zWxv0P1.png" class="tooltip" title="Solución Agrandadora"/>`,
	SolucionEncogedora:              `<img src="https://i.imgur.com/nO9UMLH.png" class="tooltip" title="Solución Encongedora"/>`,
	Veritaserum:                     `<img src="https://i.imgur.com/29AFHVA.png" class="tooltip" title="Veritaserum"/>`,
	ZumoDeMandragora:                `<img src="https://i.imgur.com/la020XU.jpg" class="tooltip" title="Zumo de Mandrágora"/>`,
}

var PotionIngredients = map[string]string{
	AgudizadoraDeIngenio:            "Escarabajos machacados\nBilis de armadillo\nRaíz de jengibre cortada",
	AlientoDeFuego:                  "Menta\nValeriana\nSemilla de fuego\nCuerno de dragón en polvo\nLavanda",
	Amortentia:                      "Asfódelo cortado\nTisana\nSemillas de anís verde\nRaíz de Angélica\nComino\nHinojo\nAcónito\nAjenjo",
	AntidotoParaVenenosComunes:      "Rocío de luna\nEsporas de vainilla de viento\nMoco de gusarajo\nAcónito\nPiel de serpiente arbórea africana\nAguamiel\nMenta\nMandrágora cocida\nEsencia de lavanda",
	AntidotoDeGlumbumble:            "Melaza de Glumbumble",
	AntidotoParaVenenosToxicos:      "Semillas de fuego\nAguijones de Billywig\nCuerno de graphorn en polvo\nCaparazones de chizpurfle",
	BrebajeParlanchin:               "Ramitas de valeriana\nAcónito\nDíctamo",
	Doxycida:                        "5 Medidas de Bundimun\n1 Hígado de dragón\n3 Caparazones de streeler\n5 Medidas de esencia de cicuta virosa\n3 Medidas de esencia de cicuta\n3 Medidas de tintura de potentilla",
	ElixirParaInducirEuforia:        "Ramitas de menta\nHigos secos\nPúas de puercoespín\nAjenjo\nSemillas de ricino\nGranos de sopóforo",
	Embellecedora:                   "Alas de hada\nRocío de la mañana\nPétalos de Rosa perfecta\n4 Medidas de pie de león cortado\nMechón de pelo de unicornio\nRaíz de jengibre",
	EsenciaDeDictamo:                "Saliva de Salamandra\nDictamo (Planta)",
	EsenciaDeMurtlap:                "Tensa\nEncurtido de tentáculos de murtlap\nMuerdago\nValeriana",
	FelixFelicis:                    "Huevo de Ashwinder\nDrimia maritima\nTentáculo de Murtlap\nTomillo\nCáscara de huevo de Occamy\nRuda\nRábano de caballo",
	FertilizanteDeEstiercolDeDragon: "Cerebro de perezoso\nCaballitos de mar voladores\nEstiércol de dragón\nMandrágora cocida\nRaíces de margarita\nBazos de rata\nTórax de libélula",
	FiltroDeMuertosEnVida:           "Ajenjo\nAsfódelo\nRaíces de valeriana\nJugo de 12 Granos de sopóforo\nCerebro de perezoso",
	FiltroDePaz:                     "Jarabe de Eléboro\nOpalo/piedra lunar/piedra de deseo\nCuerno de unicornio (polvo)\nPuas de puercoespin (polvo)",
	GraspOfDeath:                    "500gr de semillas de granada\n250 gr de Aceite de díctamo blanco (inflamable)  \n10gr de Polvillo de alas de Hada\n200 gr Jarabe de eléboro\n1 pieza de corteza del árbol vitalizante\n500gr de Rocío de la Mañana\nLágrima de un unicornio",
	Multijugos:                      "Sanguijuelas\nCrisopos\nDescurainia sophia\nCentinodia\nPolvo de cuerno de bicornio\nPiel de serpiente arbórea africana\nAlgo de la persona en la que se vaya a convertir",
	PocionCalmante:                  "Lavanda\nCorazón de cocodrilo\nMenta",
	PocionCrecepelo:                 "Cola de rata\nPúas de puercoespín\nAguijones de billywig",
	PocionCuradoraDeFurunculos:      "4 babosas cornudas\n2 púas de puercoespín\n6 colmillos de serpiente\nCebollas malolientes\nMoco de gusarajo\nRaíz de jengibre\nEspinas de shrake\nOrtiga seca",
	PocionDeDespertares:             "6 colmillos de serpiente\n4 medidas de ingrediente estándar\n6 aguijones de billywig secos\n2 ramitas de acónito",
	PocionDeErumpent:                "Cuerno de Erumpent\nPolvo de cuerno de unicornio\nEstómago de Gusano Cornudo",
	PocionDeLaMemoria:               "1 Pluma de jobberknoll\n3 Galanthus nivalis\n2 Medidas de mandrágora cocida\n2 Medidas de salvia\nHojas de alihotsy\nMenta\nOjos de anguila",
	PocionDeLaRisa:                  "Agua de manantial\nHojas de alihotsy\nAlas de billywig\n3 Púas de knarl\nPelo de puffskein\nPolvo de rábano picante\nRisa",
	PocionDelOlvido:                 "2 Gotas de agua del río Lethe\n2 Ramitas de valeriana\n2 Medidas de Ingrediente estándar\n4 Bayas de muérdago",
	PocionHerbicida:                 "Ortigas secas\nPúa de puercoespín\nColmillos de serpiente",
	PocionHerbovitalizante:          "Corteza de azarollo\nMoly\nDíctamo\nUna pinta de zumo de horklump\n2 Gotas de moco de gusarajo\n7 Colmillos de chizpurfle\nBaba de aguijón de billywig\nUna ramita de menta\nZumo de bayaboom\nMandrágora cocida\nGotas de aguamiel\nMucosa de cerebro de perezoso\nGotas de rocío de luna\nAsfódelo\nCuerno de unicornio\nAcónito\nSangre de salamandra\n10 Espinas de pez león",
	PocionOculus:                    "Ajenjo\nCuerno de unicornio\nPolvo de ópalo\nMandrágora cocida",
	PocionParaElDolorDeEstomago:     "Menta\nJarabe de Menta\nAnís\nPétalos de Manzanilla\nHígado de dragón",
	PocionPimentonica:               "Cuerno de bicornio\nRaíz de mandrágora\nImpatiens capensis",
	PocionProtectoraContraLasLlamas: "Hongo explosivo\nSangre de salamandra\nPolvos verrugosos",
	PocionVigorizante:               "Hojas de Alihotsy\nAguijones secos de billywig\nMenta\nMandrágora cocida\nInfusión de Ajenjo\nAguamiel\nInfusión de verbena\nCoclearia\nLigústico",
	PocionVolubilis:                 "Aguamiel\nRamitas de menta\nMandrágora cocida\nJarabe de eléboro",
	PocimaParaDormir:                "4 ramitas de Lavanda\n6 medidas del ingrediente estándar\n2 cucharadas de moco de gusarajo\n4 ramitas de valeriana",
	SolucionAgrandadora:             "3 Ojos de pez globo\n1 Bazo de murciélago\n2 Cucharadas de ortigas secas",
	SolucionEncogedora:              "Raíces de Margarita\nHigos secos\nOrugas lanudas\nBazo de rata\nSanguijuelas\nCicuta virosa\nAjenjo",
	Veritaserum:                     "Un pelo de cola de Unicornio adulto macho\nPluma de Fénix\nMedio litro de agua del Río Nilo (Egipto)\nUn trozo de dedo de un Grindylow\nCorazón de dragón\nAcónito\nJarabe de eléboro",
	ZumoDeMandragora:                "Mandragora",
}

var PotionNames = []string{
	AgudizadoraDeIngenio,
	AlientoDeFuego,
	Amortentia,
	AntidotoParaVenenosComunes,
	AntidotoDeGlumbumble,
	AntidotoParaVenenosToxicos,
	BrebajeParlanchin,
	Doxycida,
	ElixirParaInducirEuforia,
	Embellecedora,
	EsenciaDeDictamo,
	EsenciaDeMurtlap,
	FelixFelicis,
	FertilizanteDeEstiercolDeDragon,
	FiltroDeMuertosEnVida,
	FiltroDePaz,
	GraspOfDeath,
	Multijugos,
	PocionCalmante,
	PocionCrecepelo,
	PocionCuradoraDeFurunculos,
	PocionDeDespertares,
	PocionDeErumpent,
	PocionDeLaMemoria,
	PocionDeLaRisa,
	PocionDelOlvido,
	PocionHerbicida,
	PocionHerbovitalizante,
	PocionOculus,
	PocionParaElDolorDeEstomago,
	PocionPimentonica,
	PocionProtectoraContraLasLlamas,
	PocionVigorizante,
	PocionVolubilis,
	PocimaParaDormir,
	SolucionAgrandadora,
	SolucionEncogedora,
	Veritaserum,
	ZumoDeMandragora,
}
