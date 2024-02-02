package constants

import "strings"

// LanguageCode 语言码
//
// 采用ISO 3166-1标准
// 例如：
//
// 美式英文: en-US
//
// 简体中文: zh-CN
//
// 繁体中文: zh-TW
//
// @see https://en.wikipedia.org/wiki/ISO_3166-1
type LanguageCode string

const (
	LanguageAfrikaans                 LanguageCode = "af"
	LanguageAlbanian                  LanguageCode = "sq"
	LanguageArabicAlgeria             LanguageCode = "ar-DZ"
	LanguageArabicBahrain             LanguageCode = "ar-BH"
	LanguageArabicEgypt               LanguageCode = "ar-EG"
	LanguageArabicIraq                LanguageCode = "ar-IQ"
	LanguageArabicJordan              LanguageCode = "ar-JO"
	LanguageArabicKuwait              LanguageCode = "ar-KW"
	LanguageArabicLebanon             LanguageCode = "ar-LB"
	LanguageArabicLibya               LanguageCode = "ar-LY"
	LanguageArabicMorocco             LanguageCode = "ar-MA"
	LanguageArabicOman                LanguageCode = "ar-OM"
	LanguageArabicQatar               LanguageCode = "ar-QA"
	LanguageArabicSaudiArabia         LanguageCode = "ar-SA"
	LanguageArabicSyria               LanguageCode = "ar-SY"
	LanguageArabicTunisia             LanguageCode = "ar-TN"
	LanguageArabicUAE                 LanguageCode = "ar-AE"
	LanguageArabicYemen               LanguageCode = "ar-YE"
	LanguageBasque                    LanguageCode = "eu"
	LanguageBelarusian                LanguageCode = "be"
	LanguageBulgarian                 LanguageCode = "bg"
	LanguageCatalan                   LanguageCode = "ca"
	LanguageChineseHongKong           LanguageCode = "zh-HK"
	LanguageChinesePRC                LanguageCode = "zh-CN"
	LanguageChineseSingapore          LanguageCode = "zh-SG"
	LanguageChineseTaiwan             LanguageCode = "zh-TW"
	LanguageCroatian                  LanguageCode = "hr"
	LanguageCzech                     LanguageCode = "cs"
	LanguageDanish                    LanguageCode = "da"
	LanguageDutchBelgium              LanguageCode = "nl-BE"
	LanguageDutchStandard             LanguageCode = "nl"
	LanguageEnglish                   LanguageCode = "en"
	LanguageEnglishAustralia          LanguageCode = "en-AU"
	LanguageEnglishBelize             LanguageCode = "en-BZ"
	LanguageEnglishCanada             LanguageCode = "en-CA"
	LanguageEnglishIreland            LanguageCode = "en-IE"
	LanguageEnglishJamaica            LanguageCode = "en-JM"
	LanguageEnglishNewZealand         LanguageCode = "en-NZ"
	LanguageEnglishSouthAfrica        LanguageCode = "en-ZA"
	LanguageEnglishTrinidad           LanguageCode = "en-TT"
	LanguageEnglishUnitedKingdom      LanguageCode = "en-GB"
	LanguageEnglishUnitedStates       LanguageCode = "en-US"
	LanguageEstonian                  LanguageCode = "et"
	LanguageFaeroese                  LanguageCode = "fo"
	LanguageFarsi                     LanguageCode = "fa"
	LanguageFinnish                   LanguageCode = "fi"
	LanguageFrenchBelgium             LanguageCode = "fr-BE"
	LanguageFrenchCanada              LanguageCode = "fr-CA"
	LanguageFrenchLuxembourg          LanguageCode = "fr-LU"
	LanguageFrenchStandard            LanguageCode = "fr"
	LanguageFrenchSwitzerland         LanguageCode = "fr-CH"
	LanguageGaelicScotland            LanguageCode = "gd"
	LanguageGermanAustria             LanguageCode = "de-AT"
	LanguageGermanLiechtenstein       LanguageCode = "de-LI"
	LanguageGermanLuxembourg          LanguageCode = "de-LU"
	LanguageGermanStandard            LanguageCode = "de"
	LanguageGermanSwitzerland         LanguageCode = "de-CH"
	LanguageGreek                     LanguageCode = "el"
	LanguageHebrew                    LanguageCode = "he"
	LanguageHindi                     LanguageCode = "hi"
	LanguageHungarian                 LanguageCode = "hu"
	LanguageIcelandic                 LanguageCode = "is"
	LanguageIndonesian                LanguageCode = "id"
	LanguageIrish                     LanguageCode = "ga"
	LanguageItalianStandard           LanguageCode = "it"
	LanguageItalianSwitzerland        LanguageCode = "it-CH"
	LanguageJapanese                  LanguageCode = "ja"
	LanguageKorean                    LanguageCode = "ko"
	LanguageKoreanJohab               LanguageCode = "ko"
	LanguageKurdish                   LanguageCode = "ku"
	LanguageLatvian                   LanguageCode = "lv"
	LanguageLithuanian                LanguageCode = "lt"
	LanguageMacedonianFYROM           LanguageCode = "mk"
	LanguageMalayalam                 LanguageCode = "ml"
	LanguageMalaysian                 LanguageCode = "ms"
	LanguageMaltese                   LanguageCode = "mt"
	LanguageNorwegian                 LanguageCode = "no"
	LanguageNorwegianBokmal           LanguageCode = "nb"
	LanguageNorwegianNynorsk          LanguageCode = "nn"
	LanguagePolish                    LanguageCode = "pl"
	LanguagePortugueseBrazil          LanguageCode = "pt-BR"
	LanguagePortuguesePortugal        LanguageCode = "pt"
	LanguagePunjabi                   LanguageCode = "pa"
	LanguageRhaetoRomanic             LanguageCode = "rm"
	LanguageRomanian                  LanguageCode = "ro"
	LanguageRomanianRepublicOfMoldova LanguageCode = "ro-MD"
	LanguageRussian                   LanguageCode = "ru"
	LanguageRussianRepublicOfMoldova  LanguageCode = "ru-MD"
	LanguageSerbian                   LanguageCode = "sr"
	LanguageSlovak                    LanguageCode = "sk"
	LanguageSlovenian                 LanguageCode = "sl"
	LanguageSorbian                   LanguageCode = "sb"
	LanguageSpanishArgentina          LanguageCode = "es-AR"
	LanguageSpanishBolivia            LanguageCode = "es-BO"
	LanguageSpanishChile              LanguageCode = "es-CL"
	LanguageSpanishColombia           LanguageCode = "es-CO"
	LanguageSpanishCostaRica          LanguageCode = "es-CR"
	LanguageSpanishDominicanRepublic  LanguageCode = "es-DO"
	LanguageSpanishEcuador            LanguageCode = "es-EC"
	LanguageSpanishElSalvador         LanguageCode = "es-SV"
	LanguageSpanishGuatemala          LanguageCode = "es-GT"
	LanguageSpanishHonduras           LanguageCode = "es-HN"
	LanguageSpanishMexico             LanguageCode = "es-MX"
	LanguageSpanishNicaragua          LanguageCode = "es-NI"
	LanguageSpanishPanama             LanguageCode = "es-PA"
	LanguageSpanishParaguay           LanguageCode = "es-PY"
	LanguageSpanishPeru               LanguageCode = "es-PE"
	LanguageSpanishPuertoRico         LanguageCode = "es-PR"
	LanguageSpanishSpain              LanguageCode = "es"
	LanguageSpanishUruguay            LanguageCode = "es-UY"
	LanguageSpanishVenezuela          LanguageCode = "es-VE"
	LanguageSwedish                   LanguageCode = "sv"
	LanguageSwedishFinland            LanguageCode = "sv-FI"
	LanguageThai                      LanguageCode = "th"
	LanguageTsonga                    LanguageCode = "ts"
	LanguageTswana                    LanguageCode = "tn"
	LanguageTurkish                   LanguageCode = "tr"
	LanguageUkrainian                 LanguageCode = "ua"
	LanguageUrdu                      LanguageCode = "ur"
	LanguageVenda                     LanguageCode = "ve"
	LanguageVietnamese                LanguageCode = "vi"
	LanguageWelsh                     LanguageCode = "cy"
	LanguageXhosa                     LanguageCode = "xh"
	LanguageYiddish                   LanguageCode = "ji"
	LanguageZulu                      LanguageCode = "zu"
)

var i18n = []LanguageCode{
	LanguageAfrikaans,
	LanguageAlbanian,
	LanguageArabicAlgeria,
	LanguageArabicBahrain,
	LanguageArabicEgypt,
	LanguageArabicIraq,
	LanguageArabicJordan,
	LanguageArabicKuwait,
	LanguageArabicLebanon,
	LanguageArabicLibya,
	LanguageArabicMorocco,
	LanguageArabicOman,
	LanguageArabicQatar,
	LanguageArabicSaudiArabia,
	LanguageArabicSyria,
	LanguageArabicTunisia,
	LanguageArabicUAE,
	LanguageArabicYemen,
	LanguageBasque,
	LanguageBelarusian,
	LanguageBulgarian,
	LanguageCatalan,
	LanguageChineseHongKong,
	LanguageChinesePRC,
	LanguageChineseSingapore,
	LanguageChineseTaiwan,
	LanguageCroatian,
	LanguageCzech,
	LanguageDanish,
	LanguageDutchBelgium,
	LanguageDutchStandard,
	LanguageEnglish,
	LanguageEnglishAustralia,
	LanguageEnglishBelize,
	LanguageEnglishCanada,
	LanguageEnglishIreland,
	LanguageEnglishJamaica,
	LanguageEnglishNewZealand,
	LanguageEnglishSouthAfrica,
	LanguageEnglishTrinidad,
	LanguageEnglishUnitedKingdom,
	LanguageEnglishUnitedStates,
	LanguageEstonian,
	LanguageFaeroese,
	LanguageFarsi,
	LanguageFinnish,
	LanguageFrenchBelgium,
	LanguageFrenchCanada,
	LanguageFrenchLuxembourg,
	LanguageFrenchStandard,
	LanguageFrenchSwitzerland,
	LanguageGaelicScotland,
	LanguageGermanAustria,
	LanguageGermanLiechtenstein,
	LanguageGermanLuxembourg,
	LanguageGermanStandard,
	LanguageGermanSwitzerland,
	LanguageGreek,
	LanguageHebrew,
	LanguageHindi,
	LanguageHungarian,
	LanguageIcelandic,
	LanguageIndonesian,
	LanguageIrish,
	LanguageItalianStandard,
	LanguageItalianSwitzerland,
	LanguageJapanese,
	LanguageKorean,
	LanguageKoreanJohab,
	LanguageKurdish,
	LanguageLatvian,
	LanguageLithuanian,
	LanguageMacedonianFYROM,
	LanguageMalayalam,
	LanguageMalaysian,
	LanguageMaltese,
	LanguageNorwegian,
	LanguageNorwegianBokmal,
	LanguageNorwegianNynorsk,
	LanguagePolish,
	LanguagePortugueseBrazil,
	LanguagePortuguesePortugal,
	LanguagePunjabi,
	LanguageRhaetoRomanic,
	LanguageRomanian,
	LanguageRomanianRepublicOfMoldova,
	LanguageRussian,
	LanguageRussianRepublicOfMoldova,
	LanguageSerbian,
	LanguageSlovak,
	LanguageSlovenian,
	LanguageSorbian,
	LanguageSpanishArgentina,
	LanguageSpanishBolivia,
	LanguageSpanishChile,
	LanguageSpanishColombia,
	LanguageSpanishCostaRica,
	LanguageSpanishDominicanRepublic,
	LanguageSpanishEcuador,
	LanguageSpanishElSalvador,
	LanguageSpanishGuatemala,
	LanguageSpanishHonduras,
	LanguageSpanishMexico,
	LanguageSpanishNicaragua,
	LanguageSpanishPanama,
	LanguageSpanishParaguay,
	LanguageSpanishPeru,
	LanguageSpanishPuertoRico,
	LanguageSpanishSpain,
	LanguageSpanishUruguay,
	LanguageSpanishVenezuela,
	LanguageSwedish,
	LanguageSwedishFinland,
	LanguageThai,
	LanguageTsonga,
	LanguageTswana,
	LanguageTurkish,
	LanguageUkrainian,
	LanguageUrdu,
	LanguageVenda,
	LanguageVietnamese,
	LanguageWelsh,
	LanguageXhosa,
	LanguageYiddish,
	LanguageZulu,
}

func (lc LanguageCode) ToLowerCase() string {
	return strings.ToLower(string(lc))
}

func (lc LanguageCode) ToUpperCase() string {
	return strings.ToUpper(string(lc))
}

func (lc LanguageCode) Exist() bool {
	var exist bool
	for _, v := range i18n {
		if v == lc {
			exist = true
			break
		}
	}

	return exist
}
