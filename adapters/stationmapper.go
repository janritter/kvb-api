package adapters

import (
	"context"
	"errors"

	"github.com/sahilm/fuzzy"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

type StationMapperAdapter struct{}

func NewStationMapperAdapter() *StationMapperAdapter {
	return &StationMapperAdapter{}
}

func (adapter *StationMapperAdapter) GetStationIDForName(ctx context.Context, name string) (int, error) {
	ctx, span := otel.Tracer("kvb-api").Start(ctx, "GetStationIDForName")
	defer span.End()

	span.SetAttributes(attribute.String("input_name", name))

	foundStationName, err := findClosestMatchingStation(ctx, name)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return -1, err
	}

	return getStationIDForName(ctx, foundStationName), nil
}

// This will be replaced by a new implementation, for now this is just copied from the old code
func findClosestMatchingStation(ctx context.Context, name string) (string, error) {
	_, span := otel.Tracer("kvb-api").Start(ctx, "findClosestMatchingStation")
	defer span.End()

	wordsToTest := []string{"Aachener Str./Gürtel",
		"Adolf-Menzel-Str.",
		"Adrian-Meller-Str.",
		"Aeltgen-Dünwald-Str.",
		"Akazienweg",
		"Albin-Köbis-Straße",
		"Albrecht-Dürer-Platz",
		"Alfred-Schütte-Allee",
		"Alfter / Alanus Hochschule",
		"Alte Forststr.",
		"Alte Post",
		"Alte Römerstr.",
		"Alter Deutzer Postweg",
		"Alter Flughafen Butzweilerhof",
		"Alter Militärring",
		"Altonaer Platz",
		"Alzeyer Str.",
		"Am Bilderstöckchen",
		"Am Braunsacker",
		"Am Coloneum",
		"Am Eifeltor",
		"Am Emberg",
		"Am Faulbach",
		"Am Feldrain",
		"Am Feldrain (Sürth)",
		"Am Flachsrosterweg",
		"Am Grauen Stein",
		"Am Heiligenhäuschen",
		"Am Hetzepetsch",
		"Am Hochkreuz",
		"Am Kreuzweg",
		"Am Kölnberg", "Am Leinacker",
		"Am Lindenweg",
		"Am Neuen Forst",
		"Am Nordpark",
		"Am Portzenacker",
		"Am Schildchen",
		"Am Serviesberg",
		"Am Springborn",
		"Am Steinneuerhof",
		"Am Vorgebirgstor",
		"Am Weißen Mönch",
		"Am Zehnthof",
		"Amselstr.",
		"Amsterdamer Str./Gürtel",
		"An den Kaulen",
		"An der alten Post",
		"An der Ronne",
		"An St. Marien",
		"Andreaskloster",
		"Anemonenweg",
		"Antoniusstr.",
		"Appellhofplatz", "Arenzhof",
		"Arnoldshöhe", "Arnulfstr.",
		"Arthur-Hantzsch-Str.",
		"Auf dem Streitacker",
		"Auf der Aue",
		"Auf der Freiheit",
		"August-Horch-Str.",
		"Auguste-Kowalski-Str.",
		"Äußere Kanalstr.",
		"Autobahn",
		"Auweiler",
		"Auweilerweg",
		"Bachemer Str.",
		"Bachstelzenweg",
		"Bad Godesb. Bf/Löbestr.",
		"Bad Godesb. Bahnhof/Löbestr.",
		"Badorf",
		"Bahnstr.",
		"Baldurstr.",
		"Baptiststr.",
		"Barbarastr.",
		"Barbarossaplatz", "Baumschulenweg",
		"Bayenthalgürtel", "Beethovenstr.",
		"Belvederestr.",
		"Bensberg",
		"Bergheim Friedhof",
		"Bergheim Fährhaus",
		"Bergheim Grundschule",
		"Bergheim Industriegebiet",
		"Bergheim Kirche",
		"Bergstr.",
		"Bernkasteler Str.", "Berrenrather Str.",
		"Berrenrather Str./Gürtel",
		"Bertha-Benz-Karree",
		"Betzdorfer Str.", "Beuelsweg",
		"Beuelsweg Nord",
		"Beuthener Str.",
		"Bevingsweg",
		"Bf Deutz/LANXESS arena",
		"Bahnhof Deutz/LANXESS arena",
		"Bf Deutz/Messe",
		"Bahnhof Deutz/Messe",
		"Bf Deutz/Messeplatz",
		"Bahnhof Deutz/Messeplatz",
		"Bf Ehrenfeld",
		"Bahnhof Ehrenfeld",
		"Bf Lövenich",
		"Bahnhof Lövenich",
		"Bf Mülheim",
		"Bahnhof Mülheim",
		"Bf Porz",
		"Bahnhof Porz",
		"Bieselweg",
		"Birkenallee",
		"Birkenweg",
		"Birkenweg Schleife",
		"Bismarckstr.", "Bistritzer Str.",
		"Bitterstr.",
		"Blaugasse",
		"Blériotstr.",
		"Blockstr.",
		"Blumenberg S-Bahn",
		"Bocklemünd",
		"Bodinusstr.",
		"Boltensternstr.",
		"Bonhoefferstr.",
		"Bonn Bad Godesberg Stadthalle",
		"Bonn Bertha-von-Suttnerplatz",
		"Bonn Hauptbahnhof",
		"Bonn West",
		"Bonner Landstr.", "Bonner Str.",
		"Bonner Str./Gürtel", "Bonner Wall", "Bonntor",
		"Bornheim",
		"Bornheim Rathaus",
		"Borsigstr.",
		"Brahmsstr.",
		"Braugasse",
		"Bremerhavener Str.",
		"Breslauer Platz/Hbf",
		"Breslauer Platz/Hauptbahnhof",
		"Broichstr.",
		"Bruder-Klaus-Siedlung",
		"Brück Mauspfad",
		"Brüggener Str.", "Brühl Mitte",
		"Brühl Nord",
		"Brühl Süd",
		"Brühl-Vochem",
		"Brühler Str./Gürtel", "Brühler Straße",
		"Btf. Merheim",
		"Buchforst S-Bahn",
		"Buchforst Waldecker Str.",
		"Buchheim Frankfurter Str.",
		"Buchheim Herler Str.",
		"Buchheimer Weg",
		"Bugenhagenstr.",
		"Bundesrechnungshof",
		"Bunsenstr.",
		"Burgwiesenstr.",
		"Buschdorf",
		"Buschfeldstr.",
		"Buschweg",
		"Butzweilerstr.",
		"Bückebergstr.",
		"Böcklinstr.",
		"Bödinger Str.", "Carl-Goerdeler-Str.",
		"Carlswerkstraße",
		"Celsiusstr.",
		"Chempark S-Bahn",
		"Cheruskerstr.",
		"Chlodwigplatz", "Chorbuschstr.",
		"Chorweiler",
		"Christian-Sünner-Straße",
		"Christophstr./Mediapark", "Clarenbachstift",
		"Colonia-Allee",
		"Corintostraße",
		"Cranachstr.",
		"Curt-Stenvert-Bogen",
		"Cäsarstr.", "CöllnParc",
		"Dasselstr./Bf Süd", "Deckstein",
		"Dasselstr./Bahnhof Süd", "Deckstein",
		"Dellbrück Hauptstr.",
		"Dellbrück Mauspfad",
		"Dellbrück S-Bahn",
		"Dersdorf",
		"Deutz Technische Hochschule", "Deutzer Freiheit", "Deutzer Friedhof",
		"Deutzer Ring",
		"Diepenbeekallee",
		"Diepeschrather Str.",
		"Dieselstr.",
		"Dionysstr.",
		"DLR",
		"Dohmengasse",
		"Dom/Hbf",
		"Dom/Hauptbahnhof",
		"Donatusstr.",
		"Dornstr.",
		"Dorotheenstraße",
		"Dr.-Schultz-Str.",
		"Dransdorf",
		"Drehbrücke", "Drosselweg",
		"Dünnwald Waldbad",
		"Dünnwalder Str.",
		"Dürener Str./Gürtel",
		"Dädalusring",
		"Ebernburgweg",
		"Ebertplatz", "Ebertplatz/Riehler Str.",
		"Eddaweg",
		"Edelhofstr.",
		"Edmund-Rumpler-Str.",
		"Edsel-Ford-Str.",
		"Efferen",
		"Egelspfad",
		"Eggerbachstr.",
		"Egonstr.",
		"Eichenstr.",
		"Eifelplatz", "Eifelstr.", "Eifelwall", "Eil Heumarer Str.",
		"Eil Kirche",
		"Eiler Str.",
		"Elisabeth-Breuer-Str.",
		"Elisabethstr.",
		"Elsdorf",
		"Emilstr.",
		"Engeldorfer Hof", "Engeldorfer Str.", "Ensen Gilgaustr.",
		"Ensen Kloster",
		"Erker Mühle",
		"Erlenweg",
		"Ernst-Volland-Str.",
		"Esch",
		"Esch Friedhof",
		"Escher See",
		"Escher Str.",
		"Eschmar Bergheimer Str.",
		"Eschmar Kirche",
		"Esserstr.",
		"Esso",
		"Ettore-Bugatti-Straße",
		"Etzelstr.",
		"Eupener Str.",
		"Europaring",
		"Euskirchener Str.",
		"Eythstr.",
		"Falkenweg",
		"Feldbergstr.",
		"Feldkasseler Weg",
		"Feltenstr.",
		"Feuerwache",
		"Fischenich",
		"Flachsweg",
		"Flehbachstr.",
		"Flittard Süd",
		"Flittarder Feld",
		"Florastr.",
		"Florenzer Str.",
		"Flughafen Personalparkplatz",
		"Fordwerke Mitte",
		"Fordwerke Nord",
		"Fordwerke Süd",
		"Frankenforst",
		"Frankenstr.", "Frankenthaler Str.",
		"Frankfurter Str. S-Bahn",
		"Frankstr.",
		"Franziska-Anneke-Str.",
		"Frechen Bf",
		"Frechen Bahnhof",
		"Frechen Kirche",
		"Frechen Rathaus",
		"Frechen-Benzelrath",
		"Frechener Weg",
		"Freiheitsring",
		"Freiligrathstr.",
		"Friedenspark",
		"Friedensstr.",
		"Friedhof Chorweiler",
		"Friedhof Godorf",
		"Friedhof Lehmbacher Weg",
		"Friedhof Stammheim",
		"Friedhof Steinneuerhof",
		"Friedhof Worringen",
		"Friedrich-Hirsch-Str.",
		"Friedrich-Karl-Str./Neusser Str.",
		"Friedrich-Karl-Str./Niehler Str.",
		"Friesenplatz", "Fuldaer Str.",
		"Further Str.",
		"Fühlingen",
		"Fühlinger Weg",
		"Gaedestr.", "Gauweg",
		"Geestemünder Str.",
		"Geibelstr.",
		"Geisselstr.",
		"Geldernstr./Parkgürtel",
		"Gerhart-Hauptmann-Str.",
		"Gewerbegebiet Broichstr.",
		"Gewerbegebiet Pesch",
		"Gewerbegebiet Pesch Nord",
		"Gießener Str.",
		"Gisbertstr.",
		"Glashüttenstr.",
		"Gleueler Str./Gürtel",
		"Godorf Bf",
		"Godorf Bahnhof",
		"Goldammerweg",
		"Goldregenweg",
		"Goltsteinstr./Gürtel", "Gottesweg", "Grachtenhofstr.",
		"Graditzer Str.",
		"Graf-Adolf-Str.",
		"Gremberg",
		"Grengel Mauspfad",
		"Grevenbroicher Str.",
		"Grimmelshausenstr.",
		"Gronauer Str.",
		"Grunerstr.",
		"Grüner Weg",
		"Grüngürtelstr.",
		"Grünstr.",
		"Gummersbacher Straße",
		"Gunther-Plüschow-Str.",
		"Guntherstr.",
		"Gut Leidenhausen",
		"Gut Neuenhof",
		"Gutenbergstr.",
		"Gürzenichstr.", "Güterverkehrszentrum",
		"Güterverkehrszentrum Süd",
		"Görlinger Zentrum",
		"Göttinger Str.",
		"Habichtstraße",
		"Hackhauser Weg",
		"Hagenstr.",
		"Hahnwald",
		"Hahnwald Im Hasengarten",
		"Hahnwaldweg",
		"Halfengasse",
		"Hammerschmidtstr.",
		"Hans-Böckler-Platz/Bf West",
		"Hans-Böckler-Platz/Bahnhof West",
		"Hans-Offermann-Str.",
		"Hansaring", "Hansestr.",
		"Hansestr. Ost",
		"Hansestr. Süd",
		"Hansestr. West",
		"Haus Fühlingen",
		"Haus Vorst",
		"Havelstr.",
		"Heeresamt", "Heimersdorf",
		"Heimfriedweg",
		"Heinering",
		"Heinrich-Bützler-Straße",
		"Heinrich-Lübke-Ufer",
		"Heinrich-Mann-Str.",
		"Heinrich-Steinmann-Str.",
		"Heinz-Kühn-Str.",
		"Herforder Str.",
		"Hermann-Löns-Str.",
		"Herrigergasse",
		"Hersel",
		"Herstattallee",
		"Herthastr.", "Heumarkt", "Heussallee/Museumsmeile",
		"Hildegardis-Krankenhaus",
		"Hildegundweg",
		"Hochkirchen", "Hochkreuz",
		"Hohenlind",
		"Holweide S-Bahn",
		"Holweide Vischeringstr.",
		"Honschaftsstr.",
		"Hopfenstr.",
		"Hugo-Eckener-Str.",
		"Hugo-Junkers-Str.",
		"Humboldtstr.",
		"Hücheln Krankenhaus",
		"Hüchelner Str.",
		"Hürth Kalscheuren Bf",
		"Hürth Kalscheuren Bahnhof",
		"Hürth-Hermülheim",
		"Häuschensweg",
		"Höhenberg Frankfurter Str.",
		"Höhscheider Weg",
		"Höningen Rondorfer Weg", "Höningen Siedlung", "IKEA Am Butzweilerhof",
		"IKEA Godorf",
		"Iltisstr.",
		"Im Buschfelde",
		"Im Falkenhorst",
		"Im Hoppenkamp",
		"Im Klarenpesch",
		"Im Langen Bruch",
		"Im Rheinpark", "Im Rheintal",
		"Im Wasserfeld",
		"Im Weidenbruch",
		"Im Wichemshof",
		"Im Wirtskamp",
		"Imbacher Weg",
		"Imbuschstr.",
		"Immendorf",
		"Immendorf Schule",
		"Immendorf Siedlung",
		"Indianapolis-Straße",
		"Innere Kanalstr.",
		"Jasminweg",
		"Johannes-Prassel-Str.",
		"Johannesstr.",
		"Josef-Lammerting-Allee", "Josephstr.",
		"Junkersdorf",
		"Juridicum",
		"Justizzentrum",
		"Kalk Kapelle",
		"Kalk Post",
		"Kalk-Karree",
		"Kalker Friedhof",
		"Kalkweg",
		"Kallbergstr.",
		"Kalscheurer Weg", "Kapellenweg",
		"Kapfenberger Str.",
		"Karl-Marx-Allee",
		"Karl-Schwering-Platz",
		"Karnevalsmuseum",
		"Kartäuserhof",
		"Kaserne Haupttor",
		"Kaserne Nordtor",
		"Kasselberg",
		"Katharinenhof",
		"Kendenicher Str.", "Kesselsgasse",
		"Kettelerstr.", "Keupstr.",
		"Kiebitzweg",
		"Kieler Str.",
		"Kierberger Str.",
		"Kinderkrankenhaus",
		"Kippekausen",
		"Kirschbaumweg",
		"Kitschburger Str.",
		"Klaprothstr.",
		"Kleinfeldchensweg",
		"Kleingartenanlage Ostheim",
		"Klettenbergpark",
		"Klingerstr.",
		"Klinikum Merheim",
		"Klosterhof",
		"Koblenzer Str.", "Kochwiesenstr.",
		"Koelnmesse", "Kolkrabenweg",
		"Konrad-Adenauer-Str.",
		"Konradstr.",
		"Kopernikusschule",
		"Koppensteinstr.",
		"Kornblumenweg",
		"Krefelder Wall", "Kretzerstr.",
		"Krieger-Straße",
		"Krieler Str.",
		"Kuenstr.",
		"Kühzällerweg",
		"Kürtenstr.",
		"Kämpchensweg",
		"Köln/Bonn Flughafen",
		"Kölner Str.",
		"Kölner Weg",
		"Kölnstr.",
		"Königsforst",
		"Körnerstr.",
		"Lacher Broch",
		"Lahnstr.",
		"Langel Fähre",
		"Langel Kuhlenweg",
		"Langel Mohlenweg",
		"Langel Nord",
		"Leiblplatz",
		"Leichweg", "Leimbachweg",
		"Leinsamenweg",
		"Leipziger Platz",
		"Lenauplatz",
		"Lentpark",
		"Leopold-Gmelin-Str.",
		"Lerchenweg", "Lessingstr.",
		"Leuchterstr.",
		"Leyboldstr.", "Leyendeckerstr.",
		"Liblarer Str.", "Libur Kirche",
		"Libur Margaretenstr.",
		"Liebigstr.",
		"Lina-Bommer-Weg",
		"Lindenburg",
		"Lindenbuschweg",
		"Lindenweg",
		"Linder Kreuz",
		"Linder Mauspfad",
		"Linder Weg",
		"Lindweilerfeld",
		"Lindweilerweg",
		"Lippeweg",
		"Lohsestr.",
		"Longerich Friedhof",
		"Longerich S-Bahn",
		"Longericher Str.",
		"Longericher Str. Nord",
		"Longericher Str./Etzelstr.",
		"Lucasstr.",
		"Ludwig-Quidde-Platz",
		"Ludwigsburger Str.",
		"Lustheide",
		"LVR-Klinik",
		"Lüderichstr.",
		"Lülsdorf Hallenbad",
		"Lülsdorf Kirche",
		"Lülsdorf Nord",
		"Lülsdorf Schulzentrum",
		"Lülsdorf Uhlandstr.",
		"Maarhäuser Weg",
		"Maarweg",
		"Mannesmannstr.",
		"Mannsfeld", "Marconistr.",
		"Marconistr. Ost",
		"Margaretastr.",
		"Maria-Himmelfahrt-Str.",
		"Marienberger Weg",
		"Marienburg Südpark", "Marienburger Str.", "Marienplatz",
		"Marienstr.",
		"Marktplatz Sürth",
		"Marktstr.", "Marsdorf",
		"Maternusplatz",
		"Mathias-Brüggen-Str.",
		"Mauritiuskirche", "Mauritiusschule",
		"Max-Löbner-Str./Friesdorf",
		"Mechternstr.",
		"Meerfeldstr.",
		"Melaten",
		"Melli-Beese-Str.",
		"Mennweg",
		"Merheim",
		"Merheimer Platz",
		"Merianstr.",
		"Merkenich",
		"Merkenich Mitte",
		"Merkenicher Str.",
		"Merten",
		"Meschenich Kirche", "Messe Omnibushof",
		"Methweg",
		"Metternicherstr.",
		"Michaelshoven",
		"Militärringstr.",
		"Mohnweg",
		"Mollwitzstr.",
		"Moltkestr.", "Mommsenstr.",
		"Mondorf Ahrstr.",
		"Mondorf Beckergasse",
		"Mondorf Provinzialstr.",
		"Mondorf Rosenthalstr.",
		"Mondorf Sportplatz",
		"Montanusstr.",
		"Morsestr.",
		"Moses-Hess-Str.",
		"Mozartstr.",
		"Museum Koenig",
		"Mutzbach",
		"Mühlengasse",
		"Mühlenweg",
		"Mühlenweiher",
		"Mülhauser Str.",
		"Mülheim Berliner Str.",
		"Mülheim Wiener Platz",
		"Mülheimer Friedhof",
		"Mülheimer Ring",
		"Müllekoven",
		"Müngersdorf S-Bahn/Technologiepark",
		"Nachtigallenstr.",
		"Nattermannallee",
		"Neißestr.",
		"Nesselrodestr.",
		"Neuenweg",
		"Neuer Mülheimer Friedhof",
		"Neufelder Str.",
		"Neufeldweg",
		"Neumarkt",
		"Neurather Weg",
		"Neusser Str./Gürtel",
		"Neven DuMont Haus",
		"Nibelungenplatz",
		"Nibelungenstr.",
		"Niederkassel Evgl. Kirche",
		"Niederkassel Nord",
		"Niederkassel Rathausplatz",
		"Niederkassel Spicher Str.",
		"Niederkassel Waldstr.",
		"Niehl",
		"Niehl Betriebshof Nord",
		"Niehl Sebastianstr.",
		"Niehler Damm",
		"Niehler Kirchweg",
		"Niehler Str.",
		"Nievenheimer Str.",
		"Nippes S-Bahn",
		"Nordfriedhof",
		"Nordstr.",
		"Nußbaumerstr.",
		"Nüssenberger Str.",
		"Oberer Komarweg", "Oberlar Landgrafenstr",
		"Oberlar Lindlaustr.",
		"Oberzündorf",
		"Odenthaler Str.",
		"Ollenhauerring",
		"Ollenhauerstraße",
		"Olof-Palme-Allee",
		"Olpener Str.",
		"Oranienstr.",
		"Oranjehofstr.",
		"Oskar-Jäger-Str.",
		"Oskar-Jäger-Str./Gürtel",
		"Oskar-Schindler-Str.",
		"Ossendorf",
		"Osterather Str.",
		"Ostfriedhof",
		"Ostheim",
		"Ostlandstr.",
		"Ostmerheimer Str.",
		"Otto-Hahn-Str.",
		"Otto-Müller-Str.",
		"Palmenhof",
		"Pasteurstr.",
		"Paul-Nießen-Str.",
		"Paul-Reifenberg-Str.",
		"Pesch Schulstr.",
		"Pescher Weg",
		"Pettenkoferstr.",
		"Pierstr.",
		"Piusstr.",
		"Plittersdorfer Straße",
		"Pohligstr.", "Poll Hauptstr.",
		"Poll Salmstr.",
		"Poller Holzweg",
		"Poller Kirchweg", "Porz Markt",
		"Porz Steinstr.",
		"Porz-Langel Kirche",
		"Porz-Langel Mühle",
		"Porz-Langel Nord",
		"Porz-Langel Süd",
		"Porz-Langel Zur Eiche",
		"Porzer Str.",
		"Poststr.", "Propsthof Nord", "Prälat-van-Acken-Str.",
		"Pulheimer Str.",
		"Raiffeisenstr.",
		"Ramersdorf",
		"Ramrather Weg",
		"Ranzel Gewerbegebiet",
		"Ranzel Kirche",
		"Ranzel Schule", "Ranzel Schulstr.",
		"Ranzel Sonnenbergerweg",
		"Ranzel Weilerhof",
		"Rath-Heumar",
		"Rathaus", "Rathenaustr.",
		"Refrath",
		"Reichenspergerplatz", "Reiherstr.",
		"Reischplatz",
		"Rektor-Klein-Str.",
		"Remscheider Str.",
		"Rheidt Bahnhofstr.",
		"Rheidt Markt",
		"Rheidt Nord",
		"Rheidt Süd",
		"Rheidt Unterführung",
		"Rheinauhafen",
		"Rheinbergstr.",
		"Rheinenergie-Stadion",
		"Rheinkassel",
		"Rheinlandstr.",
		"Rheinsteinstr.", "Rhöndorfer Str.",
		"Richard-Wagner-Str.",
		"Riehler Gürtel",
		"Ritterstr.",
		"Robert-Bosch-Str.",
		"Robert-Kirchhoff-Straße",
		"Robert-Perthel-Str.",
		"Robert-Schuman-Platz",
		"Rodenkirchen Bf",
		"Rodenkirchen Bahnhof",
		"Rodenkirchen Bismarckstr.",
		"Rodenkirchen Rathaus",
		"Rodenkirchener Str.",
		"Roggenweg",
		"Roisdorf West",
		"Roisdorfer Str.", "Rolandstr.", "Rolshover Str.",
		"Rondorf", "Roonstr.", "Rosenhügel",
		"Rosenstr.", "Rosmarinweg",
		"Rotdornweg",
		"Roteichenweg",
		"Rudolf-Diesel-Str.",
		"Rudolfplatz", "Rösrather Str.",
		"Röttgensweg",
		"Saarbrücker Str.",
		"Saarstr.",
		"Sachsenbergstr.",
		"Sauerlandstr.",
		"Schadowstr.",
		"Schaffrathsgasse",
		"Schanzenstr. Nord",
		"Schanzenstr./Schauspielhaus",
		"Scheibenstr.",
		"Scheuermühlenstr.",
		"Schillingsrotter Str.",
		"Schirmerstr.",
		"Schlagbaumsweg",
		"Schlebusch",
		"Schlehdornstr.",
		"Schlettstadter Str.",
		"Schloss Röttgen",
		"Schmiedegasse",
		"Schneider-Clauss-Str.",
		"Schokoladenmuseum",
		"Schulzentrum Wahn",
		"Schumacherring",
		"Schwabenstr.",
		"Schwadorf",
		"Schwarzrheindorf Kirche",
		"Schwarzrheindorf Schule",
		"Schwarzrheindorf Siegaue",
		"Schwindstr.",
		"Schüttewerk",
		"Schützenhofstr.",
		"Schönhauser Str.", "Sechzigstr.",
		"Seeberg",
		"Seithümerstr.",
		"Selma-Lagerlöf-Str.",
		"Seniorenzentrum Riehl",
		"Servatiusstr.",
		"Severinsbrücke", "Severinskirche", "Severinstr.", "Severinusstr.",
		"Siebengebirgsallee",
		"Siedlung Mielenforst",
		"Siegburg Bf",
		"Siegburg Bahnhof",
		"Siegburg Brückberg",
		"Siegburg Ernststr.",
		"Siegburg Friedrich-Ebert-Str.",
		"Siegburg Heinrichstr.",
		"Siegburg Kaiserstr.",
		"Siegburg Kaserne",
		"Siegburg Markt",
		"Siegburg Stadthalle",
		"Siegburg Waldstr.",
		"Siegburg Zum Hohen Ufer",
		"Siegburger Str.",
		"Siegfriedstr.",
		"Sieglar Feuerwache",
		"Sieglar Flachtenstr./Krankenhaus",
		"Sieglar Im Kirschtal",
		"Sieglar Leostr.",
		"Sieglar Rathausstr.",
		"Sieglar Rathausstr./Kreisel",
		"Sieglar RSVG",
		"Sieglar Schulzentrum",
		"Siegstr.",
		"Siemensstr.",
		"Sigwinstr.",
		"Silbermöwenweg", "Sinnersdorf Kirche",
		"Sinnersdorfer Mühle",
		"Slabystr.",
		"Sparkasse",
		"Sparkasse Am Butzweilerhof",
		"Spitzangerweg",
		"Sportplatzstr.",
		"Sprengelstr.",
		"St. Vincenz Haus",
		"St. Vinzenz-Hospital",
		"St.-Tönnis-Str.",
		"St.Joseph-Kirche",
		"Stallagsweg",
		"Stammheim S-Bahn",
		"Stammheimer Ring",
		"Stegerwaldsiedlung",
		"Steinkauzweg",
		"Steinmetzstr.",
		"Steinstr. S-Bahn",
		"Steinweg",
		"Sterrenhofweg",
		"Stiftsstr.",
		"Stolberger Str.",
		"Stolberger Str./Eupener Str.",
		"Stolberger Str./Maarweg",
		"Stommeler Str.",
		"Stormstr.",
		"Straßburger Platz",
		"Stresemannstr.",
		"Stüttgenhof",
		"Stüttgerhofweg",
		"Subbelrather Str./Gürtel",
		"Suevenstr.", "Südallee",
		"Südbahnhof",
		"Sülz Hermeskeiler Platz",
		"Sülzburgstr.",
		"Sülzburgstr./Berrenrather Str.",
		"Sülzgürtel",
		"Sürth Bf",
		"Sürth Bahnhof",
		"Tacitusstr.", "Takustr.",
		"Talweg",
		"Tannenbusch Mitte",
		"Tannenbusch Süd",
		"Taubenholzweg",
		"Technologiepark Köln",
		"TechnologiePark Mitte",
		"Theodor-Heuss-Str.",
		"Theresienstr.",
		"Thermalbad",
		"Thielenbruch",
		"Thurner Kamp",
		"Trifelsstr.",
		"Trimbornstr.",
		"Troisdorf Aggerbrücke",
		"Troisdorf Altenforst",
		"Troisdorf Bergeracker",
		"Troisdorf BF",
		"Troisdorf Bahnhof",
		"Troisdorf Elsenplatz",
		"Troisdorf Kuttgasse",
		"Troisdorf Rathaus",
		"Troisdorf Ursulaplatz",
		"Troisdorf Wilhelmstr.",
		"Troisdorfer Str.",
		"TÜV-Akademie",
		"Ubierring", "Uedorf",
		"Uferstr.",
		"Ulrepforte", "Universitaet/Markt",
		"Universität",
		"Universitätsstr.",
		"Unnauer Weg",
		"Urbach Breslauer Str.",
		"Urbach Friedhof",
		"Urbach Kaiserstr.",
		"Urbach Waldstr.",
		"Urfeld",
		"Venloer Str./Gürtel",
		"Vingst",
		"Vitalisstr. Nord",
		"Vitalisstr. Süd",
		"Vogelsanger Markt",
		"Vogelsanger Str.",
		"Vogelsanger Str./Maarweg",
		"Vogelsanger Weg",
		"Volkhovener Weg",
		"Volksgarten",
		"Voltastr.",
		"Von-Galen-Str.",
		"Von-Hünefeld-Str.",
		"Von-Lohe-Str.",
		"Von-Quadt-Str.",
		"Von-Sparr-Str.",
		"Wahn Friedhof",
		"Wahn Kirche",
		"Wahn S-Bahn",
		"Waidmarkt", "Walberberg",
		"Waldorf",
		"Waldstr.",
		"Waldstr./Akazienweg",
		"Walter-Dodde-Weg",
		"Walter-Pauli-Ring",
		"Wasserwerk", "Wattstr.",
		"WDR",
		"Weichselring",
		"Weiden Einkaufszentrum",
		"Weiden Goethestr.",
		"Weiden Schulstr.",
		"Weiden Sportplatz",
		"Weiden West S-Bahn",
		"Weiden Zentrum",
		"Weidenpescher Str.",
		"Weilburger Str.",
		"Weiler",
		"Weilerweg",
		"Weinsbergstr./Gürtel",
		"Weiß Friedhof",
		"Weißer Hauptstr.",
		"Weißhausstr.",
		"Welserstr.",
		"Wendelinstr.",
		"Weserpromenade",
		"Wesseling",
		"Wesseling Nord",
		"Wesseling Süd",
		"Wesselinger Str.",
		"Westerwaldstr.", "Westfriedhof",
		"Westhoven Berliner Str.",
		"Westhoven Kölner Str.",
		"Weyertal",
		"Wezelostr.",
		"Wichheimer Str.",
		"Widdersdorf",
		"Widdersdorfer Str.",
		"Widdig",
		"Wiedenfelder Weg",
		"Wiedstr.",
		"Wiehler Str.",
		"Wiener Weg",
		"Wiesenweg",
		"Wildpark",
		"Wilhelm-Leuschner-Str.",
		"Wilhelm-Sollmann-Str.",
		"Wilhelmstr.",
		"Willi-Lauf-Allee",
		"Windmühlenstr.",
		"Wingertsheide",
		"Wiso-Fakultät",
		"Wolffsohnstr.",
		"Worringen S-Bahn",
		"Worringen Süd",
		"Worringer Str.", "Wupperplatz",
		"Wurzerstraße",
		"Wüllnerstr.",
		"Würzburger Str.",
		"Xantener Str.",
		"Zaunhof",
		"Zaunstr.",
		"Zollstock Südfriedhof", "Zollstockgürtel", "Zollstocksweg", "Zonser Str.",
		"Zoo/Flora",
		"Zugweg",
		"Zum Hedelsberg",
		"Zum Neuen Kreuz",
		"Zur Abtei",
		"Zülpicher Platz", "Zülpicher Str./Gürtel",
		"Zündorf",
		"Zündorf Altersheim",
		"Zündorf Kirche",
		"Zündorf Marktstr.",
		"Zündorf Mitte",
		"Zündorf Olefsgasse",
		"Zündorf Ranzeler Str.",
		"Zündorfer Weg",
		"Zypressenstr."}

	span.SetAttributes(attribute.String("input_name", name))

	matches := fuzzy.Find(name, wordsToTest)

	//Check if a match was found
	if len(matches) >= 1 {
		stationName := matches[0].Str
		span.SetAttributes(attribute.String("found_name", stationName))
		return stationName, nil
	}

	err := errors.New("No station for given name found")
	span.RecordError(err)
	span.SetStatus(codes.Error, err.Error())

	return "", err

}

// This will be replaced by a new implementation, for now this is just copied from the old code
func getStationIDForName(ctx context.Context, name string) int {
	_, span := otel.Tracer("kvb-api").Start(ctx, "getStationIDForName")
	defer span.End()

	var m = make(map[string]int)

	m["Aachem[ner Str./Gürtel"] = 178
	m["Adolfm[-Menzel-Str."] = 119
	m["Adrian-Meller-Str."] = 232
	m["Aeltgen-Dünwald-Str."] = 630
	m["Akazienweg"] = 264
	m["Albin-Köbis-Straße"] = 453
	m["Albrecht-Dürer-Platz"] = 755
	m["Alfred-Schütte-Allee"] = 441
	m["Alfter / Alanus Hochschule"] = 681
	m["Alte Forststr."] = 560
	m["Alte Post"] = 3665
	m["Alte Römerstr."] = 424
	m["Alter Deutzer Postweg"] = 263
	m["Alter Flughafen Butzweilerhof"] = 274
	m["Alter Militärring"] = 186
	m["Altonaer Platz"] = 363
	m["Alzeyer Str."] = 331
	m["Am Bilderstöckchen"] = 332
	m["Am Braunsacker"] = 775
	m["Am Coloneum"] = 899
	m["Am Eifeltor"] = 930
	m["Am Emberg"] = 607
	m["Am Faulbach"] = 636
	m["Am Feldrain"] = 649
	m["Am Feldrain (Sürth)"] = 905
	m["Am Flachsrosterweg"] = 942
	m["Am Grauen Stein"] = 527
	m["Am Heiligenhäuschen"] = 473
	m["Am Hetzepetsch"] = 387
	m["Am Hochkreuz"] = 898
	m["Am Kreuzweg"] = 7206
	m["Am Kölnberg"] = 85
	m["Am Leinacker"] = 877
	m["Am Lindenweg"] = 191
	m["Am Neuen Forst"] = 123
	m["Am Nordpark"] = 841
	m["Am Portzenacker"] = 619
	m["Am Schildchen"] = 787
	m["Am Serviesberg"] = 852
	m["Am Springborn"] = 611
	m["Am Steinneuerhof"] = 101
	m["Am Vorgebirgstor"] = 874
	m["Am Weißen Mönch"] = 627
	m["Am Zehnthof"] = 705
	m["Amselstr."] = 851
	m["Amsterdamer Str./Gürtel"] = 317
	m["An den Kaulen"] = 431
	m["An der alten Post"] = 889
	m["An der Ronne"] = 216
	m["An St. Marien"] = 487
	m["Andreaskloster"] = 436
	m["Anemonenweg"] = 603
	m["Antoniusstr."] = 474
	m["Appellhofplatz"] = 7
	m["Arenzhof"] = 810
	m["Arnoldshöhe"] = 84
	m["Arnulfstr."] = 151
	m["Arthur-Hantzsch-Str."] = 812
	m["Auf dem Streitacker"] = 832
	m["Auf der Aue"] = 624
	m["Auf der Freiheit"] = 754
	m["August-Horch-Str."] = 788
	m["Auguste-Kowalski-Str."] = 566
	m["Äußere Kanalstr."] = 262
	m["Autobahn"] = 534
	m["Auweiler"] = 374
	m["Auweilerweg"] = 297
	m["Bachemer Str."] = 866
	m["Bachstelzenweg"] = 287
	m["Bad Godesb. Bf/Löbestr."] = 161
	m["Bad Godesb. Bahnhof/Löbestr."] = 161
	m["Badorf"] = 737
	m["Bahnstr."] = 830
	m["Baldurstr."] = 561
	m["Baptiststr."] = 422
	m["Barbarastr."] = 646
	m["Barbarossaplatz"] = 23
	m["Baumschulenweg"] = 439
	m["Bayenthalgürtel"] = 76
	m["Beethovenstr."] = 206
	m["Belvederestr."] = 910
	m["Bensberg"] = 665
	m["Bergheim Friedhof"] = 2019
	m["Bergheim Fährhaus"] = 1069
	m["Bergheim Grundschule"] = 1010
	m["Bergheim Industriegebiet"] = 1011
	m["Bergheim Kirche"] = 2018
	m["Bergstr."] = 310
	m["Bernkasteler Str."] = 63
	m["Berrenrather Str."] = 485
	m["Berrenrather Str./Gürtel"] = 162
	m["Bertha-Benz-Karree"] = 939
	m["Betzdorfer Str."] = 48
	m["Beuelsweg"] = 827
	m["Beuelsweg Nord"] = 784
	m["Beuthener Str."] = 582
	m["Bevingsweg"] = 542
	m["Bf Deutz/LANXESS arena"] = 49
	m["Bf Deutz/Messe"] = 41
	m["Bf Deutz/Messeplatz"] = 257
	m["Bf Ehrenfeld"] = 835
	m["Bf Lövenich"] = 212
	m["Bf Mülheim"] = 572
	m["Bf Porz"] = 468
	m["Bahnhof Deutz/LANXESS arena"] = 49
	m["Bahnhof Deutz/Messe"] = 41
	m["Bahnhof Deutz/Messeplatz"] = 257
	m["Bahnhof Ehrenfeld"] = 835
	m["Bahnhof Lövenich"] = 212
	m["Bahnhof Mülheim"] = 572
	m["Bahnhof Porz"] = 468
	m["Bieselweg"] = 500
	m["Birkenallee"] = 203
	m["Birkenweg"] = 614
	m["Birkenweg Schleife"] = 803
	m["Bismarckstr."] = 33
	m["Bistritzer Str."] = 872
	m["Bitterstr."] = 429
	m["Blaugasse"] = 233
	m["Blériotstr."] = 276
	m["Blockstr."] = 399
	m["Blumenberg S-Bahn"] = 8756
	m["Bocklemünd"] = 291
	m["Bodinusstr."] = 318
	m["Boltensternstr."] = 314
	m["Bonhoefferstr."] = 642
	m["Bonn Bad Godesberg Stadthalle"] = 371
	m["Bonn Bertha-von-Suttnerplatz"] = 1115
	m["Bonn Hauptbahnhof"] = 687
	m["Bonn West"] = 688
	m["Bonner Landstr."] = 94
	m["Bonner Str."] = 457
	m["Bonner Str./Gürtel"] = 81
	m["Bonner Wall"] = 20
	m["Bonntor"] = 783
	m["Bornheim"] = 673
	m["Bornheim Rathaus"] = 628
	m["Borsigstr."] = 256
	m["Brahmsstr."] = 172
	m["Braugasse"] = 220
	m["Bremerhavener Str."] = 365
	m["Breslauer Platz/Hbf"] = 9
	m["Breslauer Platz/Hauptbahnhof"] = 9
	m["Broichstr."] = 541
	m["Bruder-Klaus-Siedlung"] = 638
	m["Brück Mauspfad"] = 547
	m["Brüggener Str."] = 62
	m["Brühl Mitte"] = 735
	m["Brühl Nord"] = 734
	m["Brühl Süd"] = 736
	m["Brühl-Vochem"] = 738
	m["Brühler Str./Gürtel"] = 74
	m["Brühler Straße"] = 689
	m["Btf. Merheim"] = 981
	m["Buchforst S-Bahn"] = 779
	m["Buchforst Waldecker Str."] = 569
	m["Buchheim Frankfurter Str."] = 577
	m["Buchheim Herler Str."] = 578
	m["Buchheimer Weg"] = 535
	m["Bugenhagenstr."] = 941
	m["Bundesrechnungshof"] = 684
	m["Bunsenstr."] = 137
	m["Burgwiesenstr."] = 591
	m["Buschdorf"] = 695
	m["Buschfeldstr."] = 590
	m["Buschweg"] = 302
	m["Butzweilerstr."] = 250
	m["Bückebergstr."] = 550
	m["Böcklinstr."] = 199
	m["Bödinger Str."] = 98
	m["Carl-Goerdeler-Str."] = 728
	m["Carlswerkstraße"] = 911
	m["Celsiusstr."] = 909
	m["Chempark S-Bahn"] = 814
	m["Cheruskerstr."] = 792
	m["Chlodwigplatz"] = 18
	m["Chorbuschstr."] = 376
	m["Chorweiler"] = 385
	m["Christian-Sünner-Straße"] = 923
	m["Christophstr./Mediapark"] = 32
	m["Clarenbachstift"] = 180
	m["Colonia-Allee"] = 592
	m["Corintostraße"] = 921
	m["Cranachstr."] = 826
	m["Curt-Stenvert-Bogen"] = 926
	m["Cäsarstr."] = 70
	m["CöllnParc"] = 900
	m["Dasselstr./Bf Süd"] = 25
	m["Dasselstr./Bahnhof Süd"] = 25
	m["Deckstein"] = 177
	m["Dellbrück Hauptstr."] = 595
	m["Dellbrück Mauspfad"] = 594
	m["Dellbrück S-Bahn"] = 604
	m["Dersdorf"] = 674
	m["Deutz Technische Hochschule"] = 44
	m["Deutzer Freiheit"] = 39
	m["Deutzer Friedhof"] = 890
	m["Deutzer Ring"] = 496
	m["Diepenbeekallee"] = 229
	m["Diepeschrather Str."] = 602
	m["Dieselstr."] = 214
	m["Dionysstr."] = 360
	m["DLR"] = 509
	m["Dohmengasse"] = 295
	m["Dom/Hb]f"] = 8
	m["Donatusstr."] = 383
	m["Dornstr."] = 430
	m["Dorotheenstraße"] = 484
	m["Dr.-Schultz-Str."] = 715
	m["Dransdorf"] = 697
	m["Drehbrücke"] = 46
	m["Drosselweg"] = 346
	m["Dünnwald Waldbad"] = 794
	m["Dünnwalder Str."] = 634
	m["Dürener Str./Gürtel"] = 170
	m["Dädalusring"] = 361
	m["Ebernburgweg"] = 330
	m["Ebertplatz"] = 35
	m["Ebertplatz/Riehler Str."] = 653
	m["Eddaweg"] = 662
	m["Edelhofstr."] = 650
	m["Edmund-Rumpler-Str."] = 867
	m["Edsel-Ford-Str."] = 414
	m["Efferen"] = 730
	m["Egelspfad"] = 211
	m["Eggerbachstr."] = 597
	m["Egonstr."] = 645
	m["Eichenstr."] = 252
	m["Eifelplatz"] = 21
	m["Eifelstr."] = 22
	m["Eifelwall"] = 26
	m["Eil Heumarer Str."] = 460
	m["Eil Kirche"] = 458
	m["Eiler Str."] = 559
	m["Elisabeth-Breuer-Str."] = 934
	m["Elisabethstr."] = 720
	m["Elsdorf"] = 481
	m["Emilstr."] = 268
	m["Engeldorfer Hof"] = 89
	m["Engeldorfer Str."] = 87
	m["Ensen Gilgaustr."] = 451
	m["Ensen Kloster"] = 452
	m["Erker Mühle"] = 548
	m["Erlenweg"] = 269
	m["Ernst-Volland-Str."] = 128
	m["Esch"] = 377
	m["Esch Friedhof"] = 857
	m["Escher See"] = 774
	m["Escher Str."] = 326
	m["Eschmar Bergheimer Str."] = 2021
	m["Eschmar Kirche"] = 2022
	m["Esserstr."] = 528
	m["Esso"] = 366
	m["Ettore-Bugatti-Straße"] = 437
	m["Etzelstr."] = 341
	m["Eupener Str."] = 181
	m["Europaring"] = 562
	m["Euskirchener Str."] = 163
	m["Eythstr."] = 514
	m["Falkenweg"] = 290
	m["Feldbergstr."] = 525
	m["Feldkasseler Weg"] = 855
	m["Feltenstr."] = 267
	m["Feuerwache"] = 471
	m["Fischenich"] = 731
	m["Flachsweg"] = 192
	m["Flehbachstr."] = 546
	m["Flittard Süd"] = 647
	m["Flittarder Feld"] = 648
	m["Florastr."] = 304
	m["Florenzer Str."] = 856
	m["Flughafen Personalparkplatz"] = 745
	m["Fordwerke Mitte"] = 369
	m["Fordwerke Nord"] = 370
	m["Fordwerke Süd"] = 368
	m["Frankenforst"] = 671
	m["Frankenstr."] = 88
	m["Frankenthaler Str."] = 329
	m["Frankfurter Str. S-Bahn"] = 657
	m["Frankstr."] = 110
	m["Franziska-Anneke-Str."] = 849
	m["Frechen Bf"] = 712
	m["Frechen Bahnhof"] = 712
	m["Frechen Kirche"] = 711
	m["Frechen Rathaus"] = 710
	m["Frechen-Benzelrath"] = 708
	m["Frechener Weg"] = 222
	m["Freiheitsring"] = 717
	m["Freiligrathstr."] = 880
	m["Friedenspark"] = 785
	m["Friedensstr."] = 477
	m["Friedhof Chorweiler"] = 398
	m["Friedhof Godorf"] = 138
	m["Friedhof Lehmbacher Weg"] = 682
	m["Friedhof Stammheim"] = 644
	m["Friedhof Steinneuerhof"] = 102
	m["Friedhof Worringen"] = 749
	m["Friedrich-Hirsch-Str."] = 480
	m["Friedrich-Karl-Str./Neusser Str."] = 838
	m["Friedrich-Karl-Str./Niehler Str."] = 345
	m["Friesenplatz"] = 30
	m["Fuldaer Str."] = 517
	m["Further Str."] = 423
	m["Fühlingen"] = 405
	m["Fühlinger Weg"] = 397
	m["Gaedestr."] = 82
	m["Gauweg"] = 821
	m["Geestemünder Str."] = 367
	m["Geibelstr."] = 146
	m["Geisselstr."] = 253
	m["Geldernstr./Parkgürtel"] = 325
	m["Gerhart-Hauptmann-Str."] = 589
	m["Gewerbegebiet Broichstr."] = 545
	m["Gewerbegebiet Pesch"] = 382
	m["Gewerbegebiet Pesch Nord"] = 781
	m["Gießener Str."] = 531
	m["Gisbertstr."] = 842
	m["Glashüttenstr."] = 862
	m["Gleueler Str./Gürtel"] = 173
	m["Godorf Bf"] = 134
	m["Godorf Bahnhof"] = 134
	m["Goldammerweg"] = 288
	m["Goldregenweg"] = 629
	m["Goltsteinstr./Gürtel"] = 78
	m["Gottesweg"] = 54
	m["Grachtenhofstr."] = 723
	m["Graditzer Str."] = 348
	m["Graf-Adolf-Str."] = 573
	m["Gremberg"] = 530
	m["Grengel Mauspfad"] = 476
	m["Grevenbroicher Str."] = 293
	m["Grimmelshausenstr."] = 112
	m["Gronauer Str."] = 580
	m["Grunerstr."] = 587
	m["Grüner Weg"] = 904
	m["Grüngürtelstr."] = 117
	m["Grünstr."] = 568
	m["Gummersbacher Straße"] = 920
	m["Gunther-Plüschow-Str."] = 661
	m["Guntherstr."] = 503
	m["Gut Leidenhausen"] = 927
	m["Gut Neuenhof"] = 722
	m["Gutenbergstr."] = 238
	m["Gürzenichstr."] = 5
	m["Güterverkehrszentrum"] = 859
	m["Güterverkehrszentrum Süd"] = 915
	m["Görlinger Zentrum"] = 300
	m["Göttinger Str."] = 823
	m["Habichtstraße"] = 884
	m["Hackhauser Weg"] = 427
	m["Hagenstr."] = 504
	m["Hahnwald"] = 105
	m["Hahnwald Im Hasengarten"] = 802
	m["Hahnwaldweg"] = 324
	m["Halfengasse"] = 350
	m["Hammerschmidtstr."] = 129
	m["Hans-Böckler-Platz/Bf West"] = 31
	m["Hans-Böckler-Platz/Bahnhof West"] = 31
	m["Hans-Offermann-Str."] = 938
	m["Hansaring"] = 36
	m["Hansestr."] = 455
	m["Hansestr. Ost"] = 454
	m["Hansestr. Süd"] = 462
	m["Hansestr. West"] = 790
	m["Haus Fühlingen"] = 404
	m["Haus Vorst"] = 713
	m["Havelstr."] = 883
	m["Heeresamt"] = 73
	m["Heimersdorf"] = 384
	m["Heimfriedweg"] = 943
	m["Heinering"] = 372
	m["Heinrich-Bützler-Straße"] = 924
	m["Heinrich-Lübke-Ufer"] = 897
	m["Heinrich-Mann-Str."] = 301
	m["Heinrich-Steinmann-Str."] = 870
	m["Heinz-Kühn-Str."] = 878
	m["Herforder Str."] = 364
	m["Hermann-Löns-Str."] = 375
	m["Herrigergasse"] = 189
	m["Hersel"] = 678
	m["Herstattallee"] = 388
	m["Herthastr."] = 53
	m["Heumark]t"] = 1
	m["Heussallee/Museumsmeile"] = 692
	m["Hildegardis-Krankenhaus"] = 145
	m["Hildegundweg"] = 618
	m["Hochkirchen"] = 95
	m["Hochkreuz"] = 699
	m["Hohenlind"] = 174
	m["Holweide S-Bahn"] = 586
	m["Holweide Vischeringstr."] = 583
	m["Honschaftsstr."] = 610
	m["Hopfenstr."] = 868
	m["Hugo-Eckener-Str."] = 275
	m["Hugo-Junkers-Str."] = 861
	m["Humboldtstr."] = 456
	m["Hücheln Krankenhaus"] = 707
	m["Hüchelner Str."] = 718
	m["Hürth Kalscheuren Bf"] = 5447
	m["Hürth Kalscheuren Bahnhof"] = 5447
	m["Hürth-Hermülheim"] = 733
	m["Häuschensweg"] = 266
	m["Höhenberg Frankfurter Str."] = 518
	m["Höhscheider Weg"] = 615
	m["Höningen Rondorfer Weg"] = 92
	m["Höningen Siedlung"] = 93
	m["IKEA Am Butzweilerhof"] = 976
	m["IKEA Godorf"] = 871
	m["Iltisstr."] = 249
	m["Im Buschfelde"] = 230
	m["Im Falkenhorst"] = 459
	m["Im Hoppenkamp"] = 655
	m["Im Klarenpesch"] = 714
	m["Im Langen Bruch"] = 552
	m["Im Rheinpark"] = 51
	m["Im Rheintal"] = 918
	m["Im Wasserfeld"] = 443
	m["Im Weidenbruch"] = 606
	m["Im Wichemshof"] = 793
	m["Im Wirtskamp"] = 795
	m["Imbacher Weg"] = 801
	m["Imbuschstr."] = 729
	m["Immendorf"] = 140
	m["Immendorf Schule"] = 840
	m["Immendorf Siedlung"] = 139
	m["Indianapolis-Straße"] = 786
	m["Innere Kanalstr."] = 756
	m["Jasminweg"] = 612
	m["Johannes-Prassel-Str."] = 378
	m["Johannesstr."] = 746
	m["Josef-Lammerting-Allee"] = 91
	m["Josephstr."] = 978
	m["Junkersdorf"] = 200
	m["Juridicum"] = 685
	m["Justizzentrum"] = 875
	m["Kalk Kapelle"] = 513
	m["Kalk Post"] = 512
	m["Kalk-Karree"] = 922
	m["Kalker Friedhof"] = 539
	m["Kalkweg"] = 622
	m["Kallbergstr."] = 928
	m["Kalscheurer Weg"] = 55
	m["Kapellenweg"] = 748
	m["Kapfenberger Str."] = 716
	m["Karl-Marx-Allee"] = 386
	m["Karl-Schwering-Platz"] = 148
	m["Karnevalsmuseum"] = 259
	m["Kartäuserhof"] = 894
	m["Kaserne Haupttor"] = 888
	m["Kaserne Nordtor"] = 482
	m["Kasselberg"] = 419
	m["Katharinenhof"] = 979
	m["Kendenicher Str."] = 60
	m["Kesselsgasse"] = 706
	m["Kettelerstr."] = 90
	m["Keupstr."] = 631
	m["Kiebitzweg"] = 732
	m["Kieler Str."] = 574
	m["Kierberger Str."] = 752
	m["Kinderkrankenhaus"] = 316
	m["Kippekausen"] = 670
	m["Kirschbaumweg"] = 103
	m["Kitschburger Str."] = 171
	m["Klaprothstr."] = 933
	m["Kleinfeldchensweg"] = 549
	m["Kleingartenanlage Ostheim"] = 321
	m["Klettenbergpark"] = 167
	m["Klingerstr."] = 863
	m["Klinikum Merheim"] = 831
	m["Klosterhof"] = 617
	m["Koblenzer Str."] = 67
	m["Kochwiesenstr."] = 296
	m["Koelnmesse]"] = 42
	m["Kolkrabenweg"] = 285
	m["Konrad-Adenauer-Str."] = 120
	m["Konradstr."] = 156
	m["Kopernikusschule"] = 470
	m["Koppensteinstr."] = 176
	m["Kornblumenweg"] = 502
	m["Krefelder Wall]"] = 38
	m["Kretzerstr."] = 309
	m["Krieger-Straße"] = 864
	m["Krieler Str."] = 175
	m["Kuenstr."] = 828
	m["Kühzällerweg"] = 588
	m["Kürtenstr."] = 522
	m["Kämpchensweg"] = 190
	m["Köln/Bonn Flughafen"] = 892
	m["Kölner Str."] = 666
	m["Kölner Weg"] = 204
	m["Kölnstr."] = 127
	m["Königsforst"] = 557
	m["Körnerstr."] = 237
	m["Lacher Broch"] = 284
	m["Lahnstr."] = 725
	m["Langel Fähre"] = 407
	m["Langel Kuhlenweg"] = 410
	m["Langel Mohlenweg"] = 409
	m["Langel Nord"] = 408
	m["Leiblplatz"] = 147
	m["Leichweg]"] = 71
	m["Leimbachweg"] = 620
	m["Leinsamenweg"] = 193
	m["Leipziger Platz"] = 307
	m["Lenauplatz"] = 247
	m["Lentpark"] = 925
	m["Leopold-Gmelin-Str."] = 817
	m["Lerchenweg"] = 96
	m["Lessingstr."] = 833
	m["Leuchterstr."] = 608
	m["Leyboldstr."] = 83
	m["Leyendeckerstr."] = 255
	m["Liblarer Str."] = 72
	m["Libur Kirche"] = 497
	m["Libur Margaretenstr."] = 498
	m["Liebigstr."] = 239
	m["Lina-Bommer-Weg"] = 937
	m["Lindenburg"] = 155
	m["Lindenbuschweg"] = 726
	m["Lindenweg"] = 847
	m["Linder Kreuz"] = 492
	m["Linder Mauspfad"] = 495
	m["Linder Weg"] = 494
	m["Lindweilerfeld"] = 393
	m["Lindweilerweg"] = 359
	m["Lippeweg"] = 616
	m["Lohsestr."] = 305
	m["Longerich Friedhof"] = 357
	m["Longerich S-Bahn"] = 358
	m["Longericher Str."] = 356
	m["Longericher Str. Nord"] = 335
	m["Longericher Str./Etzelstr."] = 860
	m["Lucasstr."] = 499
	m["Ludwig-Quidde-Platz"] = 564
	m["Ludwigsburger Str."] = 328
	m["Lustheide"] = 668
	m["LVR-Klinik"] = 843
	m["Lüderichstr."] = 919
	m["Lülsdorf Hallenbad"] = 113
	m["Lülsdorf Kirche"] = 158
	m["Lülsdorf Nord"] = 280
	m["Lülsdorf Schulzentrum"] = 281
	m["Lülsdorf Uhlandstr."] = 240
	m["Maarhäuser Weg"] = 461
	m["Maarweg"] = 179
	m["Mannesmannstr."] = 104
	m["Mannsfeld"] = 69
	m["Marconistr."] = 854
	m["Marconistr. Ost"] = 858
	m["Margaretastr."] = 273
	m["Maria-Himmelfahrt-Str."] = 584
	m["Marienberger Weg"] = 392
	m["Marienburg Südpark"] = 80
	m["Marienburger Str."] = 79
	m["Marienplatz"] = 658
	m["Marienstr."] = 834
	m["Marktplatz Sürth"] = 126
	m["Marktstr."] = 66
	m["Marsdorf"] = 235
	m["Maternusplatz"] = 109
	m["Mathias-Brüggen-Str."] = 278
	m["Mauritiuskirch]e"] = 4
	m["Mauritiusschule"] = 724
	m["Max-Löbner-Str./Friesdorf"] = 698
	m["Mechternstr."] = 771
	m["Meerfeldstr."] = 355
	m["Melaten"] = 144
	m["Melli-Beese-Str."] = 845
	m["Mennweg"] = 406
	m["Merheim"] = 540
	m["Merheimer Platz"] = 312
	m["Merianstr."] = 396
	m["Merkenich"] = 417
	m["Merkenich Mitte"] = 416
	m["Merkenicher Str."] = 349
	m["Merten"] = 676
	m["Meschenich Kirche"] = 86
	m["Messe Omnibushof"] = 501
	m["Methweg"] = 769
	m["Metternicherstr."] = 753
	m["Michaelshoven"] = 108
	m["Militärringstr."] = 279
	m["Mohnweg"] = 142
	m["Mollwitzstr."] = 336
	m["Moltkestr."] = 28
	m["Mommsenstr."] = 165
	m["Mondorf Ahrstr."] = 1071
	m["Mondorf Beckergasse"] = 1072
	m["Mondorf Provinzialstr."] = 9326
	m["Mondorf Rosenthalstr."] = 7726
	m["Mondorf Sportplatz"] = 1073
	m["Montanusstr."] = 571
	m["Morsestr."] = 853
	m["Moses-Hess-Str."] = 640
	m["Mozartstr."] = 747
	m["Museum Koenig"] = 683
	m["Mutzbach"] = 626
	m["Mühlengasse"] = 709
	m["Mühlenweg"] = 270
	m["Mühlenweiher"] = 435
	m["Mülhauser Str."] = 773
	m["Mülheim Berliner Str."] = 633
	m["Mülheim Wiener Platz"] = 570
	m["Mülheimer Friedhof"] = 519
	m["Mülheimer Ring"] = 800
	m["Müllekoven"] = 2020
	m["Müngersdorf S-Bahn/Technologiepark"] = 185
	m["Nachtigallenstr."] = 490
	m["Nattermannallee"] = 294
	m["Neißestr."] = 882
	m["Nesselrodestr."] = 818
	m["Neuenweg"] = 667
	m["Neuer Mülheimer Friedhof"] = 637
	m["Neufelder Str."] = 585
	m["Neufeldweg"] = 3729
	m["Neumarkt"] = 2
	m["Neurather Weg"] = 605
	m["Neusser Str./Gürtel"] = 303
	m["Neven DuMont Haus"] = 885
	m["Nibelungenplatz"] = 339
	m["Nibelungenstr."] = 652
	m["Niederkassel Evgl. Kirche"] = 2012
	m["Niederkassel Nord"] = 1081
	m["Niederkassel Rathausplatz"] = 2736
	m["Niederkassel Spicher Str."] = 1080
	m["Niederkassel Waldstr."] = 9327
	m["Niehl"] = 342
	m["Niehl Betriebshof Nord"] = 352
	m["Niehl Sebastianstr."] = 343
	m["Niehler Damm"] = 351
	m["Niehler Kirchweg"] = 820
	m["Niehler Str."] = 308
	m["Nievenheimer Str."] = 327
	m["Nippes S-Bahn"] = 750
	m["Nordfriedhof"] = 338
	m["Nordstr."] = 306
	m["Nußbaumerstr."] = 246
	m["Nüssenberger Str."] = 299
	m["Oberer Komarweg"] = 61
	m["Oberlar Landgrafenstr"] = 2044
	m["Oberlar Lindlaustr."] = 2043
	m["Oberzündorf"] = 763
	m["Odenthaler Str."] = 609
	m["Ollenhauerring"] = 298
	m["Ollenhauerstraße"] = 691
	m["Olof-Palme-Allee"] = 690
	m["Olpener Str."] = 551
	m["Oranienstr."] = 523
	m["Oranjehofstr."] = 412
	m["Oskar-Jäger-Str."] = 258
	m["Oskar-Jäger-Str./Gürtel"] = 182
	m["Oskar-Schindler-Str."] = 929
	m["Ossendorf"] = 271
	m["Osterather Str."] = 815
	m["Ostfriedhof"] = 598
	m["Ostheim"] = 533
	m["Ostlandstr."] = 227
	m["Ostmerheimer Str."] = 543
	m["Otto-Hahn-Str."] = 136
	m["Otto-Müller-Str."] = 381
	m["Palmenhof"] = 809
	m["Pasteurstr."] = 822
	m["Paul-Nießen-Str."] = 916
	m["Paul-Reifenberg-Str."] = 621
	m["Pesch Schulstr."] = 373
	m["Pescher Weg"] = 380
	m["Pettenkoferstr."] = 772
	m["Pierstr."] = 135
	m["Piusstr."] = 236
	m["Plittersdorfer Straße"] = 701
	m["Pohligstr."] = 52
	m["Poll Hauptstr."] = 442
	m["Poll Salmstr."] = 438
	m["Poller Holzweg"] = 445
	m["Poller Kirchweg"] = 47
	m["Porz Markt"] = 467
	m["Porz Steinstr."] = 466
	m["Porz-Langel Kirche"] = 767
	m["Porz-Langel Mühle"] = 765
	m["Porz-Langel Nord"] = 764
	m["Porz-Langel Süd"] = 703
	m["Porz-Langel Zur Eiche"] = 766
	m["Porzer Str."] = 554
	m["Poststr."] = 3
	m["Propsthof Nord"] = 43
	m["Prälat-van-Acken-Str."] = 850
	m["Pulheimer Str."] = 391
	m["Raiffeisenstr."] = 444
	m["Ramersdorf"] = 1584
	m["Ramrather Weg"] = 751
	m["Ranzel Gewerbegebiet"] = 322
	m["Ranzel Kirche"] = 344
	m["Ranzel Schule"] = 77
	m["Ranzel Schulstr."] = 6578
	m["Ranzel Sonnenbergerweg"] = 1082
	m["Ranzel Weilerhof"] = 869
	m["Rath-Heumar"] = 556
	m["Rathaus"] = 6
	m["Rathenaustr."] = 7610
	m["Refrath"] = 669
	m["Reichenspergerplatz"] = 34
	m["Reiherstr."] = 777
	m["Reischplatz"] = 798
	m["Rektor-Klein-Str."] = 272
	m["Remscheider Str."] = 516
	m["Rheidt Bahnhofstr."] = 1079
	m["Rheidt Markt"] = 1076
	m["Rheidt Nord"] = 1078
	m["Rheidt Süd"] = 1074
	m["Rheidt Unterführung"] = 1077
	m["Rheinauhafen"] = 744
	m["Rheinbergstr."] = 581
	m["Rheinenergie-Stadion"] = 187
	m["Rheinkassel"] = 411
	m["Rheinlandstr."] = 413
	m["Rheinsteinstr."] = 64
	m["Rhöndorfer Str."] = 159
	m["Richard-Wagner-Str."] = 118
	m["Riehler Gürtel"] = 319
	m["Ritterstr."] = 130
	m["Robert-Bosch-Str."] = 415
	m["Robert-Kirchhoff-Straße"] = 696
	m["Robert-Perthel-Str."] = 334
	m["Robert-Schuman-Platz"] = 1655
	m["Rodenkirchen Bf"] = 106
	m["Rodenkirchen Bahnhof"] = 106
	m["Rodenkirchen Bismarckstr."] = 780
	m["Rodenkirchen Rathaus"] = 111
	m["Rodenkirchener Str."] = 9338
	m["Roggenweg"] = 194
	m["Roisdorf West"] = 672
	m["Roisdorfer Str."] = 59
	m["Rolandstr."] = 16
	m["Rolshover Str."] = 433
	m["Rondorf"] = 97
	m["Roonstr."] = 29
	m["Rosenhügel"] = 228
	m["Rosenstr."] = 13
	m["Rosmarinweg"] = 808
	m["Rotdornweg"] = 721
	m["Roteichenweg"] = 879
	m["Rudolf-Diesel-Str."] = 463
	m["Rudolfplatz"] = 27
	m["Rösrather Str."] = 565
	m["Röttgensweg"] = 555
	m["Saarbrücker Str."] = 536
	m["Saarstr."] = 218
	m["Sachsenbergstr."] = 529
	m["Sauerlandstr."] = 532
	m["Schadowstr."] = 768
	m["Schaffrathsgasse"] = 242
	m["Schanzenstr. Nord"] = 908
	m["Schanzenstr./Schauspielhaus"] = 891
	m["Scheibenstr."] = 337
	m["Scheuermühlenstr."] = 887
	m["Schillingsrotter Str."] = 122
	m["Schirmerstr."] = 770
	m["Schlagbaumsweg"] = 593
	m["Schlebusch"] = 663
	m["Schlehdornstr."] = 931
	m["Schlettstadter Str."] = 418
	m["Schloss Röttgen"] = 558
	m["Schmiedegasse"] = 340
	m["Schneider-Clauss-Str."] = 829
	m["Schokoladenmuseum"] = 719
	m["Schulzentrum Wahn"] = 886
	m["Schumacherring"] = 895
	m["Schwabenstr."] = 121
	m["Schwadorf"] = 739
	m["Schwarzrheindorf Kirche"] = 1510
	m["Schwarzrheindorf Schule"] = 1514
	m["Schwarzrheindorf Siegaue"] = 1515
	m["Schwindstr."] = 198
	m["Schüttewerk"] = 440
	m["Schützenhofstr."] = 935
	m["Schönhauser Str."] = 65
	m["Sechzigstr."] = 311
	m["Seeberg"] = 395
	m["Seithümerstr."] = 213
	m["Selma-Lagerlöf-Str."] = 221
	m["Seniorenzentrum Riehl"] = 320
	m["Servatiusstr."] = 537
	m["Severinsbrücke"] = 45
	m["Severinskirche"] = 15
	m["Severinstr."] = 11
	m["Severinusstr."] = 219
	m["Siebengebirgsallee"] = 168
	m["Siedlung Mielenforst"] = 599
	m["Siegburg Bf"] = 1811
	m["Siegburg Bahnhof"] = 1811
	m["Siegburg Brückberg"] = 2099
	m["Siegburg Ernststr."] = 2100
	m["Siegburg Friedrich-Ebert-Str."] = 7699
	m["Siegburg Heinrichstr."] = 4969
	m["Siegburg Kaiserstr."] = 2102
	m["Siegburg Kaserne"] = 2650
	m["Siegburg Markt"] = 2105
	m["Siegburg Stadthalle"] = 2103
	m["Siegburg Waldstr."] = 2101
	m["Siegburg Zum Hohen Ufer"] = 8758
	m["Siegburger Str."] = 448
	m["Siegfriedstr."] = 116
	m["Sieglar Feuerwache"] = 2045
	m["Sieglar Flachtenstr./Krankenhaus"] = 2024
	m["Sieglar Im Kirschtal"] = 2023
	m["Sieglar Leostr."] = 2046
	m["Sieglar Rathausstr."] = 2025
	m["Sieglar Rathausstr./Kreisel"] = 2812
	m["Sieglar RSVG"] = 2026
	m["Sieglar Schulzentrum"] = 2057
	m["Siegstr."] = 107
	m["Siemensstr."] = 469
	m["Sigwinstr."] = 613
	m["Silbermöwenweg"] = 14
	m["Sinnersdorf Kirche"] = 704
	m["Sinnersdorfer Mühle"] = 379
	m["Slabystr."] = 315
	m["Sparkasse"] = 643
	m["Sparkasse Am Butzweilerhof"] = 903
	m["Spitzangerweg"] = 217
	m["Sportplatzstr."] = 656
	m["Sprengelstr."] = 819
	m["St. Vincenz Haus"] = 7394
	m["St. Vinzenz-Hospital"] = 824
	m["St.-Tönnis-Str."] = 426
	m["St.Joseph-Kirche"] = 354
	m["Stallagsweg"] = 390
	m["Stammheim S-Bahn"] = 813
	m["Stammheimer Ring"] = 641
	m["Stegerwaldsiedlung"] = 567
	m["Steinkauzweg"] = 286
	m["Steinmetzstr."] = 515
	m["Steinstr. S-Bahn"] = 625
	m["Steinweg"] = 553
	m["Sterrenhofweg"] = 205
	m["Stiftsstr."] = 9029
	m["Stolberger Str."] = 778
	m["Stolberger Str./Eupener Str."] = 183
	m["Stolberger Str./Maarweg"] = 184
	m["Stommeler Str."] = 839
	m["Stormstr."] = 225
	m["Straßburger Platz"] = 563
	m["Stresemannstr."] = 465
	m["Stüttgenhof"] = 234
	m["Stüttgerhofweg"] = 210
	m["Subbelrather Str./Gürtel"] = 245
	m["Suevenstr."] = 40
	m["Südallee"] = 207
	m["Südbahnhof"] = 876
	m["Sülz Hermeskeiler Platz"] = 166
	m["Sülzburgstr."] = 152
	m["Sülzburgstr./Berrenrather Str."] = 157
	m["Sülzgürtel"] = 160
	m["Sürth Bf"] = 124
	m["Sürth Bahnhof"] = 124
	m["Tacitusstr."] = 68
	m["Takustr."] = 836
	m["Talweg"] = 805
	m["Tannenbusch Mitte"] = 694
	m["Tannenbusch Süd"] = 693
	m["Taubenholzweg"] = 446
	m["Technologiepark Köln"] = 197
	m["TechnologiePark Mitte"] = 848
	m["Theodor-Heuss-Str."] = 464
	m["Theresienstr."] = 149
	m["Thermalbad"] = 806
	m["Thielenbruch"] = 596
	m["Thurner Kamp"] = 600
	m["Trifelsstr."] = 789
	m["Trimbornstr."] = 816
	m["Troisdorf Aggerbrücke"] = 2098
	m["Troisdorf Altenforst"] = 2083
	m["Troisdorf Bergeracker"] = 2069
	m["Troisdorf BF"] = 2071
	m["Troisdorf Bahnhof"] = 2071
	m["Troisdorf Elsenplatz"] = 2078
	m["Troisdorf Kuttgasse"] = 2664
	m["Troisdorf Rathaus"] = 2041
	m["Troisdorf Ursulaplatz"] = 2076
	m["Troisdorf Wilhelmstr."] = 2073
	m["Troisdorfer Str."] = 493
	m["TÜV-Akademie"] = 799
	m["Ubierring"] = 17
	m["Uedorf"] = 679
	m["Uferstr."] = 114
	m["Ulrepforte"] = 19
	m["Universitaet/Markt"] = 686
	m["Universität"] = 153
	m["Universitätsstr."] = 143
	m["Unnauer Weg"] = 394
	m["Urbach Breslauer Str."] = 472
	m["Urbach Friedhof"] = 479
	m["Urbach Kaiserstr."] = 511
	m["Urbach Waldstr."] = 510
	m["Urfeld"] = 743
	m["Venloer Str./Gürtel"] = 251
	m["Vingst"] = 521
	m["Vitalisstr. Nord"] = 659
	m["Vitalisstr. Süd"] = 195
	m["Vogelsanger Markt"] = 289
	m["Vogelsanger Str."] = 265
	m["Vogelsanger Str./Maarweg"] = 260
	m["Vogelsanger Weg"] = 202
	m["Volkhovener Weg"] = 402
	m["Volksgarten"] = 913
	m["Voltastr."] = 901
	m["Von-Galen-Str."] = 639
	m["Von-Hünefeld-Str."] = 277
	m["Von-Lohe-Str."] = 635
	m["Von-Quadt-Str."] = 601
	m["Von-Sparr-Str."] = 632
	m["Wahn Friedhof"] = 491
	m["Wahn Kirche"] = 489
	m["Wahn S-Bahn"] = 488
	m["Waidmarkt"] = 12
	m["Walberberg"] = 677
	m["Waldorf"] = 675
	m["Waldstr."] = 208
	m["Waldstr./Akazienweg"] = 475
	m["Walter-Dodde-Weg"] = 421
	m["Walter-Pauli-Ring"] = 243
	m["Wasserwerk"] = 75
	m["Wattstr."] = 873
	m["WDR"] = 292
	m["Weichselring"] = 403
	m["Weiden Einkaufszentrum"] = 791
	m["Weiden Goethestr."] = 226
	m["Weiden Schulstr."] = 241
	m["Weiden Sportplatz"] = 224
	m["Weiden West S-Bahn"] = 702
	m["Weiden Zentrum"] = 261
	m["Weidenpescher Str."] = 347
	m["Weilburger Str."] = 526
	m["Weiler"] = 401
	m["Weilerweg"] = 811
	m["Weinsbergstr./Gürtel"] = 254
	m["Weiß Friedhof"] = 132
	m["Weißer Hauptstr."] = 131
	m["Weißhausstr."] = 150
	m["Welserstr."] = 651
	m["Wendelinstr."] = 188
	m["Weserpromenade"] = 881
	m["Wesseling"] = 741
	m["Wesseling Nord"] = 742
	m["Wesseling Süd"] = 740
	m["Wesselinger Str."] = 125
	m["Westerwaldstr."] = 99
	m["Westfriedhof"] = 283
	m["Westhoven Berliner Str."] = 450
	m["Westhoven Kölner Str."] = 449
	m["Weyertal"] = 154
	m["Wezelostr."] = 400
	m["Wichheimer Str."] = 579
	m["Widdersdorf"] = 231
	m["Widdersdorfer Str."] = 196
	m["Widdig"] = 680
	m["Wiedenfelder Weg"] = 428
	m["Wiedstr."] = 654
	m["Wiehler Str."] = 544
	m["Wiener Weg"] = 209
	m["Wiesenweg"] = 478
	m["Wildpark"] = 623
	m["Wilhelm-Leuschner-Str."] = 727
	m["Wilhelm-Sollmann-Str."] = 362
	m["Wilhelmstr."] = 825
	m["Willi-Lauf-Allee"] = 244
	m["Windmühlenstr."] = 2548
	m["Wingertsheide"] = 7054
	m["Wiso-Fakultät"] = 846
	m["Wolffsohnstr."] = 282
	m["Worringen S-Bahn"] = 420
	m["Worringen Süd"] = 425
	m["Worringer Str."] = 37
	m["Wupperplatz"] = 944
	m["Wurzerstraße"] = 700
	m["Wüllnerstr."] = 169
	m["Würzburger Str."] = 520
	m["Xantener Str."] = 323
	m["Zaunhof"] = 141
	m["Zaunstr."] = 215
	m["Zollstock Südfriedhof"] = 57
	m["Zollstockgürtel"] = 56
	m["Zollstocksweg"] = 58
	m["Zonser Str."] = 837
	m["Zoo/Flora"] = 313
	m["Zugweg"] = 914
	m["Zum Hedelsberg"] = 133
	m["Zum Neuen Kreuz"] = 807
	m["Zur Abtei"] = 804
	m["Zülpicher Platz"] = 24
	m["Zülpicher Str./Gürtel"] = 164
	m["Zündorf"] = 486
	m["Zündorf Altersheim"] = 759
	m["Zündorf Kirche"] = 758
	m["Zündorf Marktstr."] = 757
	m["Zündorf Mitte"] = 760
	m["Zündorf Olefsgasse"] = 761
	m["Zündorf Ranzeler Str."] = 762
	m["Zündorfer Weg"] = 447
	m["Zypressenstr."] = 389

	span.SetAttributes(attribute.String("input_name", name))

	stationID := m[name]

	span.SetAttributes(attribute.Int("found_id", stationID))

	return stationID
}
