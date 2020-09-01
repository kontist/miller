package cli

import (
	"flag"
	"fmt"
	"os"

	"miller/clitypes"
	"miller/mappers"
	"miller/mapping"
)

// xxx to own file/package -- ?
var VERSION_STRING string = "v6.0.0-dev"

// ================================================================
// Miller command-line interface
// ================================================================

// xxx comment why Go flag package is used by verb handlers but not by main

//// ----------------------------------------------------------------
//#define DEFAULT_OFMT                     "%lf"
//#define DEFAULT_OQUOTING                 QUOTE_MINIMAL
//#define DEFAULT_JSON_FLATTEN_SEPARATOR   ":"
//#define DEFAULT_OOSVAR_FLATTEN_SEPARATOR ":"
//#define DEFAULT_COMMENT_STRING           "#"
//
//// ASCII 1f and 1e
//#define ASV_FS "\x1f"
//#define ASV_RS "\x1e"
//
//#define ASV_FS_FOR_HELP "0x1f"
//#define ASV_RS_FOR_HELP "0x1e"
//
//// Unicode code points U+241F and U+241E, encoded as UTF-8.
//#define USV_FS "\xe2\x90\x9f"
//#define USV_RS "\xe2\x90\x9e"
//
//#define USV_FS_FOR_HELP "U+241F (UTF-8 0xe2909f)"
//#define USV_RS_FOR_HELP "U+241E (UTF-8 0xe2909e)"

// ----------------------------------------------------------------
func ParseCommandLine(args []string) (
	options clitypes.TOptions,
	recordMappers []mapping.IRecordMapper,
	filenames []string,
	err error,
) {
	options = clitypes.DefaultOptions()
	argc := len(args)
	argi := 1

	// Try .mlrrc overrides (then command-line on top of that).
	// A --norc flag (if provided) must come before all other options.
	// Or, they can set the environment variable MLRRC="__none__".
	//	if (argc >= 2 && args[1] == "--norc" {
	//		argi++;
	//	} else {
	//		loadMlrrcOrDie(popts);
	//	}

	for argi < argc /* variable increment: 1 or 2 depending on flag */ {
		if args[argi][0] != '-' {
			break // No more flag options to process
		} else if handleTerminalUsage(args, argc, argi) {
			os.Exit(0)
		} else if handleReaderOptions(args, argc, &argi, &options.ReaderOptions) {
			// handled
		} else if handleWriterOptions(args, argc, &argi, &options.WriterOptions) {
			// handled
		} else if handleReaderWriterOptions(args, argc, &argi, &options.ReaderOptions, &options.WriterOptions) {
			// handled
		} else if handleMiscOptions(args, argc, &argi, &options) {
			// handled
		} else {
			// unhandled
			usageUnrecognizedVerb(args[0], args[argi])
		}
	}

	// xxx to do
	//	cli_apply_defaults(popts);

	//	lhmss_t* default_rses = get_default_rses();
	//	lhmss_t* default_fses = get_default_fses();
	//	lhmss_t* default_pses = get_default_pses();
	//	lhmsll_t* default_repeat_ifses = get_default_repeat_ifses();
	//	lhmsll_t* default_repeat_ipses = get_default_repeat_ipses();
	//
	//	if (options.ReaderOptions.IRS == nil)
	//		options.ReaderOptions.IRS = lhmss_get_or_die(default_rses, options.ReaderOptions.InputFileFormat);
	//	if (options.ReaderOptions.IFS == nil)
	//		options.ReaderOptions.IFS = lhmss_get_or_die(default_fses, options.ReaderOptions.InputFileFormat);
	//	if (options.ReaderOptions.ips == nil)
	//		options.ReaderOptions.ips = lhmss_get_or_die(default_pses, options.ReaderOptions.InputFileFormat);
	//
	//	if (options.ReaderOptions.allow_repeat_ifs == NEITHER_TRUE_NOR_FALSE)
	//		options.ReaderOptions.allow_repeat_ifs = lhmsll_get_or_die(default_repeat_ifses, options.ReaderOptions.InputFileFormat);
	//	if (options.ReaderOptions.allow_repeat_ips == NEITHER_TRUE_NOR_FALSE)
	//		options.ReaderOptions.allow_repeat_ips = lhmsll_get_or_die(default_repeat_ipses, options.ReaderOptions.InputFileFormat);
	//
	//	if (options.WriterOptions.ORS == nil)
	//		options.WriterOptions.ORS = lhmss_get_or_die(default_rses, options.WriterOptions.OutputFileFormat);
	//	if (options.WriterOptions.OFS == nil)
	//		options.WriterOptions.OFS = lhmss_get_or_die(default_fses, options.WriterOptions.OutputFileFormat);
	//	if (options.WriterOptions.ops == nil)
	//		options.WriterOptions.ops = lhmss_get_or_die(default_pses, options.WriterOptions.OutputFileFormat);
	//
	//	if options.WriterOptions.OutputFileFormat == "pprint") && len(options.WriterOptions.OFS) != 1) {
	//		fmt.Fprintf(os.Stderr, "%s: OFS for PPRINT format must be single-character; got \"%s\".\n",
	//			os.Args[0], options.WriterOptions.OFS);
	//		return nil;
	//	}

	//	// Construct the mapper list for single use, e.g. the normal streaming case wherein the
	//	// mappers operate on all input files. Also retain information needed to construct them
	//	// for each input file, for in-place mode.
	//	options.mapper_argb = argi;
	//	options.original_argv = args;
	//	options.non_in_place_argv = copy_argv(args);
	//	options.argc = argc;
	//	*ppmapper_list = cli_parse_mappers(options.non_in_place_argv, &argi, argc, popts);
	recordMappers, err = parseMappers(args, &argi, argc, &options)
	if err != nil {
		return options, recordMappers, filenames, err
	}

	filenames = args[argi:]

	//	if (options.no_input) {
	//		slls_free(options.filenames);
	//		options.filenames = nil;
	//	}

	//	if (options.do_in_place && (options.filenames == nil || options.filenames.length == 0)) {
	//		fmt.Fprintf(os.Stderr, "%s: -I option (in-place operation) requires input files.\n", os.Args[0]);
	//		os.Exit(1);
	//	}

	//	if (options.have_rand_seed) {
	//		mtrand_init(options.rand_seed);
	//	} else {
	//		mtrand_init_default();
	//	}

	return options, recordMappers, filenames, nil
}

// ----------------------------------------------------------------
// Returns a list of mappers, from the starting point in args given by *pargi.
// Bumps *pargi to point to remaining post-mapper-setup args, i.e. filenames.

func parseMappers(args []string, pargi *int, argc int, options *clitypes.TOptions) ([]mapping.IRecordMapper, error) {
	mapperList := make([]mapping.IRecordMapper, 0)
	argi := *pargi

	// Allow then-chains to start with an initial 'then': 'mlr verb1 then verb2 then verb3' or
	// 'mlr then verb1 then verb2 then verb3'. Particuarly useful in backslashy scripting contexts.
	if (argc-argi) >= 1 && args[argi] == "then" {
		argi++
	}

	if (argc - argi) < 1 {
		fmt.Fprintf(os.Stderr, "%s: no verb supplied.\n", os.Args[0])
		mainUsageShort()
		os.Exit(1)
	}

	for {
		checkArgCount(args, argi, argc, 1)
		verb := args[argi]

		mapperSetup := lookUpMapperSetup(verb)
		if mapperSetup == nil {
			fmt.Fprintf(os.Stderr,
				"%s: verb \"%s\" not found. Please use \"%s --help\" for a list.\n",
				os.Args[0], verb, os.Args[0])
			os.Exit(1)
		}

		// It's up to the parse func to print its usage on CLI-parse failure.
		// Also note: this assumes main reader/writer opts are all parsed
		// *before* mapper parse-CLI methods are invoked.

		mapper := mapperSetup.ParseCLIFunc(
			&argi,
			argc,
			args,
			flag.ExitOnError,
			&options.ReaderOptions,
			&options.WriterOptions,
		)

		if mapper == nil {
			// Error message already printed out
			os.Exit(1)
		}

		//		if (mapperSetup.IgnoresInput && len(mapperList) == 0) {
		//			// e.g. then-chain starts with seqgen
		//			options.no_input = true;
		//		}

		mapperList = append(mapperList, mapper)

		if argi >= argc || args[argi] != "then" {
			break
		}
		argi++
	}

	*pargi = argi
	return mapperList, nil
}

// ----------------------------------------------------------------
func lookUpMapperSetup(verb string) *mapping.MapperSetup {
	for _, mapperSetup := range MAPPER_LOOKUP_TABLE {
		if mapperSetup.Verb == verb {
			return &mapperSetup
		}
	}
	return nil
}

func listAllVerbsRaw(o *os.File) {
	for _, mapperSetup := range MAPPER_LOOKUP_TABLE {
		fmt.Fprintf(o, "%s\n", mapperSetup.Verb)
	}
}

func listAllVerbs(o *os.File, leader string) {
	separator := " "

	leaderlen := len(leader)
	separatorlen := len(separator)
	linelen := leaderlen
	j := 0

	for _, mapperSetup := range MAPPER_LOOKUP_TABLE {
		verb := mapperSetup.Verb
		verblen := len(verb)
		linelen += separatorlen + verblen
		if linelen >= 80 {
			fmt.Fprintf(o, "\n")
			linelen = leaderlen + separatorlen + verblen
			j = 0
		}
		if j == 0 {
			fmt.Fprintf(o, "%s", leader)
		}
		fmt.Fprintf(o, "%s%s", separator, verb)
		j++
	}

	fmt.Fprintf(o, "\n")
}

func usageAllVerbs(argv0 string) {
	separator := "================================================================"

	for _, mapperSetup := range MAPPER_LOOKUP_TABLE {
		fmt.Printf("%s\n", separator)
		args := [3]string{os.Args[0], mapperSetup.Verb, "--help"}
		argi := 1
		mapperSetup.ParseCLIFunc(
			&argi,
			3,
			args[:],
			flag.ContinueOnError,
			nil,
			nil,
		)
		fmt.Printf("\n")
	}
	fmt.Printf("%s\n", separator)
	os.Exit(0)
}

var MAPPER_LOOKUP_TABLE = []mapping.MapperSetup{
	mappers.CatSetup,
	mappers.NothingSetup,
	mappers.PutSetup,
	mappers.TacSetup,
}

//	&mapper_altkv_setup,
//	&mapper_bar_setup,
//	&mapper_bootstrap_setup,
//	&mapper_cat_setup,
//	&mapper_check_setup,
//	&mapper_clean_whitespace_setup,
//	&mapper_count_setup,
//	&mapper_count_distinct_setup,
//	&mapper_count_similar_setup,
//	&mapper_cut_setup,
//	&mapper_decimate_setup,
//	&mapper_fill_down_setup,
//	&mapper_filter_setup,
//	&mapper_format_values_setup,
//	&mapper_fraction_setup,
//	&mapper_grep_setup,
//	&mapper_group_by_setup,
//	&mapper_group_like_setup,
//	&mapper_having_fields_setup,
//	&mapper_head_setup,
//	&mapper_histogram_setup,
//	&mapper_join_setup,
//	&mapper_label_setup,
//	&mapper_least_frequent_setup,
//	&mapper_merge_fields_setup,
//	&mapper_most_frequent_setup,
//	&mapper_nest_setup,
//	&mapper_nothing_setup,
//	&mapper_put_setup,
//	&mapper_regularize_setup,
//	&mapper_remove_empty_columns_setup,
//	&mapper_rename_setup,
//	&mapper_reorder_setup,
//	&mapper_repeat_setup,
//	&mapper_reshape_setup,
//	&mapper_sample_setup,
//	&mapper_sec2gmt_setup,
//	&mapper_sec2gmtdate_setup,
//	&mapper_seqgen_setup,
//	&mapper_shuffle_setup,
//	&mapper_skip_trivial_records_setup,
//	&mapper_sort_setup,
//	// xxx temp for 5.4.0 -- will continue work after
//	// &mapper_sort_within_records_setup,
//	&mapper_stats1_setup,
//	&mapper_stats2_setup,
//	&mapper_step_setup,
//	&mapper_tac_setup,
//	&mapper_tail_setup,
//	&mapper_tee_setup,
//	&mapper_top_setup,
//	&mapper_uniq_setup,
//	&mapper_unsparsify_setup,

// ----------------------------------------------------------------
//static lhmss_t* singleton_pdesc_to_chars_map = nil;
//static lhmss_t* get_desc_to_chars_map() {
//	if (singleton_pdesc_to_chars_map == nil) {
//		singleton_pdesc_to_chars_map = lhmss_alloc();
//		lhmss_put(singleton_pdesc_to_chars_map, "cr",        "\r",       NO_FREE);
//		lhmss_put(singleton_pdesc_to_chars_map, "crcr",      "\r\r",     NO_FREE);
//		lhmss_put(singleton_pdesc_to_chars_map, "newline",   "\n",       NO_FREE);
//		lhmss_put(singleton_pdesc_to_chars_map, "lf",        "\n",       NO_FREE);
//		lhmss_put(singleton_pdesc_to_chars_map, "lflf",      "\n\n",     NO_FREE);
//		lhmss_put(singleton_pdesc_to_chars_map, "crlf",      "\r\n",     NO_FREE);
//		lhmss_put(singleton_pdesc_to_chars_map, "crlfcrlf",  "\r\n\r\n", NO_FREE);
//		lhmss_put(singleton_pdesc_to_chars_map, "tab",       "\t",       NO_FREE);
//		lhmss_put(singleton_pdesc_to_chars_map, "space",     " ",        NO_FREE);
//		lhmss_put(singleton_pdesc_to_chars_map, "comma",     ",",        NO_FREE);
//		lhmss_put(singleton_pdesc_to_chars_map, "newline",   "\n",       NO_FREE);
//		lhmss_put(singleton_pdesc_to_chars_map, "pipe",      "|",        NO_FREE);
//		lhmss_put(singleton_pdesc_to_chars_map, "slash",     "/",        NO_FREE);
//		lhmss_put(singleton_pdesc_to_chars_map, "colon",     ":",        NO_FREE);
//		lhmss_put(singleton_pdesc_to_chars_map, "semicolon", ";",        NO_FREE);
//		lhmss_put(singleton_pdesc_to_chars_map, "equals",    "=",        NO_FREE);
//	}
//	return singleton_pdesc_to_chars_map;
//}

func SeparatorFromArg(arg string) string {
	// xxx stub
	return arg
	//	char* chars = lhmss_get(get_desc_to_chars_map(), arg);
	//	if (chars != nil) // E.g. crlf
	//		return mlr_strdup_or_die(chars);
	//	else // E.g. '\r\n'
	//		return mlr_alloc_unbackslash(arg);
}

//// ----------------------------------------------------------------
//static lhmss_t* singleton_default_rses = nil;
//static lhmss_t* singleton_default_fses = nil;
//static lhmss_t* singleton_default_pses = nil;
//static lhmsll_t* singleton_default_repeat_ifses = nil;
//static lhmsll_t* singleton_default_repeat_ipses = nil;
//
//static lhmss_t* get_default_rses() {
//	if (singleton_default_rses == nil) {
//		singleton_default_rses = lhmss_alloc();
//
//		lhmss_put(singleton_default_rses, "gen",      "N/A",  NO_FREE);
//		lhmss_put(singleton_default_rses, "dkvp",     "auto",  NO_FREE);
//		lhmss_put(singleton_default_rses, "json",     "auto",  NO_FREE);
//		lhmss_put(singleton_default_rses, "nidx",     "auto",  NO_FREE);
//		lhmss_put(singleton_default_rses, "csv",      "auto",  NO_FREE);
//		lhmss_put(singleton_default_rses, "csvlite",  "auto",  NO_FREE);
//		lhmss_put(singleton_default_rses, "markdown", "auto",  NO_FREE);
//		lhmss_put(singleton_default_rses, "pprint",   "auto",  NO_FREE);
//		lhmss_put(singleton_default_rses, "xtab",     "(N/A)", NO_FREE);
//	}
//	return singleton_default_rses;
//}
//
//static lhmss_t* get_default_fses() {
//	if (singleton_default_fses == nil) {
//		singleton_default_fses = lhmss_alloc();
//		lhmss_put(singleton_default_fses, "gen",      "(N/A)",  NO_FREE);
//		lhmss_put(singleton_default_fses, "dkvp",     ",",      NO_FREE);
//		lhmss_put(singleton_default_fses, "json",     "(N/A)",  NO_FREE);
//		lhmss_put(singleton_default_fses, "nidx",     " ",      NO_FREE);
//		lhmss_put(singleton_default_fses, "csv",      ",",      NO_FREE);
//		lhmss_put(singleton_default_fses, "csvlite",  ",",      NO_FREE);
//		lhmss_put(singleton_default_fses, "markdown", "(N/A)",  NO_FREE);
//		lhmss_put(singleton_default_fses, "pprint",   " ",      NO_FREE);
//		lhmss_put(singleton_default_fses, "xtab",     "auto",   NO_FREE);
//	}
//	return singleton_default_fses;
//}
//
//static lhmss_t* get_default_pses() {
//	if (singleton_default_pses == nil) {
//		singleton_default_pses = lhmss_alloc();
//		lhmss_put(singleton_default_pses, "gen",      "(N/A)", NO_FREE);
//		lhmss_put(singleton_default_pses, "dkvp",     "=",     NO_FREE);
//		lhmss_put(singleton_default_pses, "json",     "(N/A)", NO_FREE);
//		lhmss_put(singleton_default_pses, "nidx",     "(N/A)", NO_FREE);
//		lhmss_put(singleton_default_pses, "csv",      "(N/A)", NO_FREE);
//		lhmss_put(singleton_default_pses, "csvlite",  "(N/A)", NO_FREE);
//		lhmss_put(singleton_default_pses, "markdown", "(N/A)", NO_FREE);
//		lhmss_put(singleton_default_pses, "pprint",   "(N/A)", NO_FREE);
//		lhmss_put(singleton_default_pses, "xtab",     " ",     NO_FREE);
//	}
//	return singleton_default_pses;
//}
//
//static lhmsll_t* get_default_repeat_ifses() {
//	if (singleton_default_repeat_ifses == nil) {
//		singleton_default_repeat_ifses = lhmsll_alloc();
//		lhmsll_put(singleton_default_repeat_ifses, "gen",      false, NO_FREE);
//		lhmsll_put(singleton_default_repeat_ifses, "dkvp",     false, NO_FREE);
//		lhmsll_put(singleton_default_repeat_ifses, "json",     false, NO_FREE);
//		lhmsll_put(singleton_default_repeat_ifses, "csv",      false, NO_FREE);
//		lhmsll_put(singleton_default_repeat_ifses, "csvlite",  false, NO_FREE);
//		lhmsll_put(singleton_default_repeat_ifses, "markdown", false, NO_FREE);
//		lhmsll_put(singleton_default_repeat_ifses, "nidx",     false, NO_FREE);
//		lhmsll_put(singleton_default_repeat_ifses, "xtab",     false, NO_FREE);
//		lhmsll_put(singleton_default_repeat_ifses, "pprint",   true,  NO_FREE);
//	}
//	return singleton_default_repeat_ifses;
//}
//
//static lhmsll_t* get_default_repeat_ipses() {
//	if (singleton_default_repeat_ipses == nil) {
//		singleton_default_repeat_ipses = lhmsll_alloc();
//		lhmsll_put(singleton_default_repeat_ipses, "gen",      false, NO_FREE);
//		lhmsll_put(singleton_default_repeat_ipses, "dkvp",     false, NO_FREE);
//		lhmsll_put(singleton_default_repeat_ipses, "json",     false, NO_FREE);
//		lhmsll_put(singleton_default_repeat_ipses, "csv",      false, NO_FREE);
//		lhmsll_put(singleton_default_repeat_ipses, "csvlite",  false, NO_FREE);
//		lhmsll_put(singleton_default_repeat_ipses, "markdown", false, NO_FREE);
//		lhmsll_put(singleton_default_repeat_ipses, "nidx",     false, NO_FREE);
//		lhmsll_put(singleton_default_repeat_ipses, "xtab",     true,  NO_FREE);
//		lhmsll_put(singleton_default_repeat_ipses, "pprint",   false, NO_FREE);
//	}
//	return singleton_default_repeat_ipses;
//}
//
//func free_opt_singletons() {
//	lhmss_free(singleton_pdesc_to_chars_map);
//	lhmss_free(singleton_default_rses);
//	lhmss_free(singleton_default_fses);
//	lhmss_free(singleton_default_pses);
//	lhmsll_free(singleton_default_repeat_ifses);
//	lhmsll_free(singleton_default_repeat_ipses);
//}
//
//// For displaying the default separators in on-line help
//static char* rebackslash(char* sep) {
//	if sep == "\r"))
//		return "\\r";
//	else if sep == "\n"))
//		return "\\n";
//	else if sep == "\r\n"))
//		return "\\r\\n";
//	else if sep == "\t"))
//		return "\\t";
//	else if sep == " "))
//		return "space";
//	else
//		return sep;
//}

// ----------------------------------------------------------------
func mainUsageShort() {
	fmt.Fprintf(os.Stderr, "Please run \"%s --help\" for detailed usage information.\n", os.Args[0])
	os.Exit(1)
}

// ----------------------------------------------------------------
// The mainUsageLong() function is split out into subroutines in support of the
// manpage autogenerator.

func mainUsageLong(o *os.File, argv0 string) {
	mainUsageSynopsis(o, argv0)
	fmt.Fprintf(o, "\n")

	fmt.Fprintf(o, "COMMAND-LINE-SYNTAX EXAMPLES:\n")
	mainUsageExamples(o, argv0, "  ")
	fmt.Fprintf(o, "\n")

	fmt.Fprintf(o, "DATA-FORMAT EXAMPLES:\n")
	mainUsageDataFormatExamples(o, argv0)
	fmt.Fprintf(o, "\n")

	fmt.Fprintf(o, "HELP OPTIONS:\n")
	mainUsageHelpOptions(o, argv0)
	fmt.Fprintf(o, "\n")

	//	fmt.Fprintf(o, "CUSTOMIZATION VIA .MLRRC:\N");
	//	mainUsageMlrrc(o, argv0);
	//	fmt.Fprintf(o, "\n");

	fmt.Fprintf(o, "VERBS:\n")
	listAllVerbs(o, "  ")
	fmt.Fprintf(o, "\n")

	//	fmt.Fprintf(o, "FUNCTIONS FOR THE FILTER AND PUT VERBS:\n");
	//	mainUsageFunctions(o, argv0, "  ");
	//	fmt.Fprintf(o, "\n");

	fmt.Fprintf(o, "DATA-FORMAT OPTIONS, FOR INPUT, OUTPUT, OR BOTH:\n")
	mainUsageDataFormatOptions(o, argv0)
	fmt.Fprintf(o, "\n")

	//	fmt.Fprintf(o, "COMMENTS IN DATA:\N");
	//	mainUsageCommentsInData(o, argv0);
	//	fmt.Fprintf(o, "\n");
	//
	fmt.Fprintf(o, "FORMAT-CONVERSION KEYSTROKE-SAVER OPTIONS:\n")
	mainUsageFormatConversionKeystrokeSaverOptions(o, argv0)
	fmt.Fprintf(o, "\n")

	//	fmt.Fprintf(o, "COMPRESSED-DATA OPTIONS:\N");
	//	mainUsageCompressedDataOptions(o, argv0);
	//	fmt.Fprintf(o, "\n");
	//
	//	fmt.Fprintf(o, "SEPARATOR OPTIONS:\n");
	//	mainUsageSeparatorOptions(o, argv0);
	//	fmt.Fprintf(o, "\n");
	//
	//	fmt.Fprintf(o, "RELEVANT TO CSV/CSV-LITE INPUT ONLY:\n");
	//	mainUsageCsvOptions(o, argv0);
	//	fmt.Fprintf(o, "\n");
	//
	//	fmt.Fprintf(o, "DOUBLE-QUOTING FOR CSV OUTPUT:\n");
	//	mainUsageDoubleQuoting(o, argv0);
	//	fmt.Fprintf(o, "\n");
	//
	//	fmt.Fprintf(o, "NUMERICAL FORMATTING:\n");
	//	mainUsageNumericalFormatting(o, argv0);
	//	fmt.Fprintf(o, "\n");
	//
	//	fmt.Fprintf(o, "OTHER OPTIONS:\n");
	//	mainUsageOtherOptions(o, argv0);
	//	fmt.Fprintf(o, "\n");
	//
	fmt.Fprintf(o, "THEN-CHAINING:\n")
	mainUsageThenChaining(o, argv0)
	fmt.Fprintf(o, "\n")

	//	fmt.Fprintf(o, "AUXILIARY COMMANDS:\N");
	//	mainUsageAuxents(o, argv0);
	//	fmt.Fprintf(o, "\n");
	//
	fmt.Fprintf(o, "SEE ALSO:\n")
	mainUsageSeeAlso(o, argv0)
}

func mainUsageSynopsis(o *os.File, argv0 string) {
	fmt.Fprintf(o, "Usage: %s [I/O options] {verb} [verb-dependent options ...] {zero or more file names}\n", argv0)
}

func mainUsageExamples(o *os.File, argv0 string, leader string) {
	fmt.Fprintf(o, "%s%s --csv cut -f hostname,uptime mydata.csv\n", leader, argv0)
	fmt.Fprintf(o, "%s%s --tsv --rs lf filter '$status != \"down\" && $upsec >= 10000' *.tsv\n", leader, argv0)
	fmt.Fprintf(o, "%s%s --nidx put '$sum = $7 < 0.0 ? 3.5 : $7 + 2.1*$8' *.dat\n", leader, argv0)
	fmt.Fprintf(o, "%sgrep -v '^#' /etc/group | %s --ifs : --nidx --opprint label group,pass,gid,member then sort -f group\n", leader, argv0)
	fmt.Fprintf(o, "%s%s join -j account_id -f accounts.dat then group-by account_name balances.dat\n", leader, argv0)
	fmt.Fprintf(o, "%s%s --json put '$attr = sub($attr, \"([0-9]+)_([0-9]+)_.*\", \"\\1:\\2\")' data/*.json\n", leader, argv0)
	fmt.Fprintf(o, "%s%s stats1 -a min,mean,max,p10,p50,p90 -f flag,u,v data/*\n", leader, argv0)
	fmt.Fprintf(o, "%s%s stats2 -a linreg-pca -f u,v -g shape data/*\n", leader, argv0)
	fmt.Fprintf(o, "%s%s put -q '@sum[$a][$b] += $x; end {emit @sum, \"a\", \"b\"}' data/*\n", leader, argv0)
	fmt.Fprintf(o, "%s%s --from estimates.tbl put '\n", leader, argv0)
	fmt.Fprintf(o, "  for (k,v in $*) {\n")
	fmt.Fprintf(o, "    if (is_numeric(v) && k =~ \"^[t-z].*$\") {\n")
	fmt.Fprintf(o, "      $sum += v; $count += 1\n")
	fmt.Fprintf(o, "    }\n")
	fmt.Fprintf(o, "  }\n")
	fmt.Fprintf(o, "  $mean = $sum / $count # no assignment if count unset'\n")
	fmt.Fprintf(o, "%s%s --from infile.dat put -f analyze.mlr\n", leader, argv0)
	fmt.Fprintf(o, "%s%s --from infile.dat put 'tee > \"./taps/data-\".$a.\"-\".$b, $*'\n", leader, argv0)
	fmt.Fprintf(o, "%s%s --from infile.dat put 'tee | \"gzip > ./taps/data-\".$a.\"-\".$b.\".gz\", $*'\n", leader, argv0)
	fmt.Fprintf(o, "%s%s --from infile.dat put -q '@v=$*; dump | \"jq .[]\"'\n", leader, argv0)
	fmt.Fprintf(o, "%s%s --from infile.dat put  '(NR %% 1000 == 0) { print > os.Stderr, \"Checkpoint \".NR}'\n",
		leader, argv0)
}

func mainUsageHelpOptions(o *os.File, argv0 string) {
	fmt.Fprintf(o, "  -h or --help                 Show this message.\n")
	fmt.Fprintf(o, "  --version                    Show the software version.\n")
	fmt.Fprintf(o, "  {verb name} --help           Show verb-specific help.\n")
	fmt.Fprintf(o, "  --help-all-verbs             Show help on all verbs.\n")
	fmt.Fprintf(o, "  -l or --list-all-verbs       List only verb names.\n")
	fmt.Fprintf(o, "  -L                           List only verb names, one per line.\n")
	fmt.Fprintf(o, "  -f or --help-all-functions   Show help on all built-in functions.\n")
	fmt.Fprintf(o, "  -F                           Show a bare listing of built-in functions by name.\n")
	fmt.Fprintf(o, "  -k or --help-all-keywords    Show help on all keywords.\n")
	fmt.Fprintf(o, "  -K                           Show a bare listing of keywords by name.\n")
}

//func mainUsageMlrrc(o *os.File, argv0 string) {
//	fmt.Fprintf(o, "You can set up personal defaults via a $HOME/.mlrrc and/or ./.mlrrc.\n");
//	fmt.Fprintf(o, "For example, if you usually process CSV, then you can put \"--csv\" in your .mlrrc file\n");
//	fmt.Fprintf(o, "and that will be the default input/output format unless otherwise specified on the command line.\n");
//	fmt.Fprintf(o, "\n");
//	fmt.Fprintf(o, "The .mlrrc file format is one \"--flag\" or \"--option value\" per line, with the leading \"--\" optional.\n");
//	fmt.Fprintf(o, "Hash-style comments and blank lines are ignored.\n");
//	fmt.Fprintf(o, "\n");
//	fmt.Fprintf(o, "Sample .mlrrc:\n");
//	fmt.Fprintf(o, "# Input and output formats are CSV by default (unless otherwise specified\n");
//	fmt.Fprintf(o, "# on the mlr command line):\n");
//	fmt.Fprintf(o, "csv\n");
//	fmt.Fprintf(o, "# These are no-ops for CSV, but when I do use JSON output, I want these\n");
//	fmt.Fprintf(o, "# pretty-printing options to be used:\n");
//	fmt.Fprintf(o, "jvstack\n");
//	fmt.Fprintf(o, "jlistwrap\n");
//	fmt.Fprintf(o, "\n");
//	fmt.Fprintf(o, "How to specify location of .mlrrc:\n");
//	fmt.Fprintf(o, "* If $MLRRC is set:\n");
//	fmt.Fprintf(o, "  o If its value is \"__none__\" then no .mlrrc files are processed.\n");
//	fmt.Fprintf(o, "  o Otherwise, its value (as a filename) is loaded and processed. If there are syntax\n");
//	fmt.Fprintf(o, "    errors, they abort mlr with a usage message (as if you had mistyped something on the\n");
//	fmt.Fprintf(o, "    command line). If the file can't be loaded at all, though, it is silently skipped.\n");
//	fmt.Fprintf(o, "  o Any .mlrrc in your home directory or current directory is ignored whenever $MLRRC is\n");
//	fmt.Fprintf(o, "    set in the environment.\n");
//	fmt.Fprintf(o, "* Otherwise:\n");
//	fmt.Fprintf(o, "  o If $HOME/.mlrrc exists, it's then processed as above.\n");
//	fmt.Fprintf(o, "  o If ./.mlrrc exists, it's then also processed as above.\n");
//	fmt.Fprintf(o, "  (I.e. current-directory .mlrrc defaults are stacked over home-directory .mlrrc defaults.)\n");
//	fmt.Fprintf(o, "\n");
//	fmt.Fprintf(o, "See also:\n");
//	fmt.Fprintf(o, "https://johnkerl.org/miller/doc/customization.html\n");
//}

//func mainUsageFunctions(o *os.File, argv0 string, char* leader) {
//	fmgr_t* pfmgr = fmgr_alloc();
//	fmgr_list_functions(pfmgr, o, leader);
//	fmgr_free(pfmgr, nil);
//	fmt.Fprintf(o, "\n");
//	fmt.Fprintf(o, "Please use \"%s --help-function {function name}\" for function-specific help.\n", argv0);
//}

func mainUsageDataFormatExamples(o *os.File, argv0 string) {
	fmt.Fprintf(o,
		`  DKVP: delimited key-value pairs (Miller default format)
  +---------------------+
  | apple=1,bat=2,cog=3 | Record 1: "apple" => "1", "bat" => "2", "cog" => "3"
  | dish=7,egg=8,flint  | Record 2: "dish" => "7", "egg" => "8", "3" => "flint"
  +---------------------+

  NIDX: implicitly numerically indexed (Unix-toolkit style)
  +---------------------+
  | the quick brown     | Record 1: "1" => "the", "2" => "quick", "3" => "brown"
  | fox jumped          | Record 2: "1" => "fox", "2" => "jumped"
  +---------------------+

  CSV/CSV-lite: comma-separated values with separate header line
  +---------------------+
  | apple,bat,cog       |
  | 1,2,3               | Record 1: "apple => "1", "bat" => "2", "cog" => "3"
  | 4,5,6               | Record 2: "apple" => "4", "bat" => "5", "cog" => "6"
  +---------------------+

  Tabular JSON: nested objects are supported, although arrays within them are not:
  +---------------------+
  | {                   |
  |  "apple": 1,        | Record 1: "apple" => "1", "bat" => "2", "cog" => "3"
  |  "bat": 2,          |
  |  "cog": 3           |
  | }                   |
  | {                   |
  |   "dish": {         | Record 2: "dish:egg" => "7", "dish:flint" => "8", "garlic" => ""
  |     "egg": 7,       |
  |     "flint": 8      |
  |   },                |
  |   "garlic": ""      |
  | }                   |
  +---------------------+

  PPRINT: pretty-printed tabular
  +---------------------+
  | apple bat cog       |
  | 1     2   3         | Record 1: "apple => "1", "bat" => "2", "cog" => "3"
  | 4     5   6         | Record 2: "apple" => "4", "bat" => "5", "cog" => "6"
  +---------------------+

  XTAB: pretty-printed transposed tabular
  +---------------------+
  | apple 1             | Record 1: "apple" => "1", "bat" => "2", "cog" => "3"
  | bat   2             |
  | cog   3             |
  |                     |
  | dish 7              | Record 2: "dish" => "7", "egg" => "8"
  | egg  8              |
  +---------------------+

  Markdown tabular (supported for output only):
  +-----------------------+
  | | apple | bat | cog | |
  | | ---   | --- | --- | |
  | | 1     | 2   | 3   | | Record 1: "apple => "1", "bat" => "2", "cog" => "3"
  | | 4     | 5   | 6   | | Record 2: "apple" => "4", "bat" => "5", "cog" => "6"
  +-----------------------+
`)
}

func mainUsageDataFormatOptions(o *os.File, argv0 string) {
	fmt.Fprintf(o,
		`
	  --idkvp   --odkvp   --dkvp      Delimited key-value pairs, e.g "a=1,b=2"
	                                  (this is Miller's default format).

	  --inidx   --onidx   --nidx      Implicitly-integer-indexed fields
	                                  (Unix-toolkit style).
	  -T                              Synonymous with "--nidx --fs tab".

	  --icsv    --ocsv    --csv       Comma-separated value (or tab-separated
	                                  with --fs tab, etc.)

	  --itsv    --otsv    --tsv       Keystroke-savers for "--icsv --ifs tab",
	                                  "--ocsv --ofs tab", "--csv --fs tab".
	  --iasv    --oasv    --asv       Similar but using ASCII FS %s and RS %s\n",
		ASV_FS_FOR_HELP, ASV_RS_FOR_HELP);
	  --iusv    --ousv    --usv       Similar but using Unicode FS %s\n",
		USV_FS_FOR_HELP);
	                                  and RS %s\n",
		USV_RS_FOR_HELP);

	  --icsvlite --ocsvlite --csvlite Comma-separated value (or tab-separated
	                                  with --fs tab, etc.). The 'lite' CSV does not handle
	                                  RFC-CSV double-quoting rules; is slightly faster;
	                                  and handles heterogeneity in the input stream via
	                                  empty newline followed by new header line. See also
	                                  http://johnkerl.org/miller/doc/file-formats.html#CSV/TSV/etc.

	  --itsvlite --otsvlite --tsvlite Keystroke-savers for "--icsvlite --ifs tab",
	                                  "--ocsvlite --ofs tab", "--csvlite --fs tab".
	  -t                              Synonymous with --tsvlite.
	  --iasvlite --oasvlite --asvlite Similar to --itsvlite et al. but using ASCII FS %s and RS %s\n",
		ASV_FS_FOR_HELP, ASV_RS_FOR_HELP);
	  --iusvlite --ousvlite --usvlite Similar to --itsvlite et al. but using Unicode FS %s\n",
		USV_FS_FOR_HELP);
	                                  and RS %s\n",
		USV_RS_FOR_HELP);

	  --ipprint --opprint --pprint    Pretty-printed tabular (produces no
	                                  output until all input is in).
	                      --right     Right-justifies all fields for PPRINT output.
	                      --barred    Prints a border around PPRINT output
	                                  (only available for output).

	            --omd                 Markdown-tabular (only available for output).

	  --ixtab   --oxtab   --xtab      Pretty-printed vertical-tabular.
	                      --xvright   Right-justifies values for XTAB format.

	  --ijson   --ojson   --json      JSON tabular: sequence or list of one-level
	                                  maps: {...}{...} or [{...},{...}].
	    --json-map-arrays-on-input    JSON arrays are unmillerable. --json-map-arrays-on-input
	    --json-skip-arrays-on-input   is the default: arrays are converted to integer-indexed
	    --json-fatal-arrays-on-input  maps. The other two options cause them to be skipped, or
	                                  to be treated as errors.  Please use the jq tool for full
	                                  JSON (pre)processing.
	                      --jvstack   Put one key-value pair per line for JSON
	                                  output.
	                --jsonx --ojsonx  Keystroke-savers for --json --jvstack
	                --jsonx --ojsonx  and --ojson --jvstack, respectively.
	                      --jlistwrap Wrap JSON output in outermost [ ].
	                    --jknquoteint Do not quote non-string map keys in JSON output.
	                     --jvquoteall Quote map values in JSON output, even if they're
	                                  numeric.
	              --jflatsep {string} Separator for flattening multi-level JSON keys,
	                                  e.g. '{"a":{"b":3}}' becomes a:b => 3 for
	                                  non-JSON formats. Defaults to %s.\n",
		DEFAULT_JSON_FLATTEN_SEPARATOR);

	  -p is a keystroke-saver for --nidx --fs space --repifs

	  Examples: --csv for CSV-formatted input and output; --idkvp --opprint for
	  DKVP-formatted input and pretty-printed output.

	  Please use --iformat1 --oformat2 rather than --format1 --oformat2.
	  The latter sets up input and output flags for format1, not all of which
	  are overridden in all cases by setting output format to format2.

`)
}

//func mainUsageCommentsInData(o *os.File, argv0 string) {
//	fmt.Fprintf(o, "  --skip-comments                 Ignore commented lines (prefixed by \"%s\")\n",
//		DEFAULT_COMMENT_STRING);
//	fmt.Fprintf(o, "                                  within the input.\n");
//	fmt.Fprintf(o, "  --skip-comments-with {string}   Ignore commented lines within input, with\n");
//	fmt.Fprintf(o, "                                  specified prefix.\n");
//	fmt.Fprintf(o, "  --pass-comments                 Immediately print commented lines (prefixed by \"%s\")\n",
//		DEFAULT_COMMENT_STRING);
//	fmt.Fprintf(o, "                                  within the input.\n");
//	fmt.Fprintf(o, "  --pass-comments-with {string}   Immediately print commented lines within input, with\n");
//	fmt.Fprintf(o, "                                  specified prefix.\n");
//	fmt.Fprintf(o, "Notes:\n");
//	fmt.Fprintf(o, "* Comments are only honored at the start of a line.\n");
//	fmt.Fprintf(o, "* In the absence of any of the above four options, comments are data like\n");
//	fmt.Fprintf(o, "  any other text.\n");
//	fmt.Fprintf(o, "* When pass-comments is used, comment lines are written to standard output\n");
//	fmt.Fprintf(o, "  immediately upon being read; they are not part of the record stream.\n");
//	fmt.Fprintf(o, "  Results may be counterintuitive. A suggestion is to place comments at the\n");
//	fmt.Fprintf(o, "  start of data files.\n");
//}

func mainUsageFormatConversionKeystrokeSaverOptions(o *os.File, argv0 string) {
	fmt.Fprintf(o, "As keystroke-savers for format-conversion you may use the following:\n")
	fmt.Fprintf(o, "        --c2t --c2d --c2n --c2j --c2x --c2p --c2m\n")
	fmt.Fprintf(o, "  --t2c       --t2d --t2n --t2j --t2x --t2p --t2m\n")
	fmt.Fprintf(o, "  --d2c --d2t       --d2n --d2j --d2x --d2p --d2m\n")
	fmt.Fprintf(o, "  --n2c --n2t --n2d       --n2j --n2x --n2p --n2m\n")
	fmt.Fprintf(o, "  --j2c --j2t --j2d --j2n       --j2x --j2p --j2m\n")
	fmt.Fprintf(o, "  --x2c --x2t --x2d --x2n --x2j       --x2p --x2m\n")
	fmt.Fprintf(o, "  --p2c --p2t --p2d --p2n --p2j --p2x       --p2m\n")
	fmt.Fprintf(o, "The letters c t d n j x p m refer to formats CSV, TSV, DKVP, NIDX, JSON, XTAB,\n")
	fmt.Fprintf(o, "PPRINT, and markdown, respectively. Note that markdown format is available for\n")
	fmt.Fprintf(o, "output only.\n")
}

//func mainUsageCompressedDataOptions(o *os.File, argv0 string) {
//	fmt.Fprintf(o, "  --prepipe {command} This allows Miller to handle compressed inputs. You can do\n");
//	fmt.Fprintf(o, "  without this for single input files, e.g. \"gunzip < myfile.csv.gz | %s ...\".\n",
//		argv0);
//	fmt.Fprintf(o, "  However, when multiple input files are present, between-file separations are\n");
//	fmt.Fprintf(o, "  lost; also, the FILENAME variable doesn't iterate. Using --prepipe you can\n");
//	fmt.Fprintf(o, "  specify an action to be taken on each input file. This pre-pipe command must\n");
//	fmt.Fprintf(o, "  be able to read from standard input; it will be invoked with\n");
//	fmt.Fprintf(o, "    {command} < {filename}.\n");
//	fmt.Fprintf(o, "  Examples:\n");
//	fmt.Fprintf(o, "    %s --prepipe 'gunzip'\n", argv0);
//	fmt.Fprintf(o, "    %s --prepipe 'zcat -cf'\n", argv0);
//	fmt.Fprintf(o, "    %s --prepipe 'xz -cd'\n", argv0);
//	fmt.Fprintf(o, "    %s --prepipe cat\n", argv0);
//	fmt.Fprintf(o, "  Note that this feature is quite general and is not limited to decompression\n");
//	fmt.Fprintf(o, "  utilities. You can use it to apply per-file filters of your choice.\n");
//	fmt.Fprintf(o, "  For output compression (or other) utilities, simply pipe the output:\n");
//	fmt.Fprintf(o, "    %s ... | {your compression command}\n", argv0);
//}

//func mainUsageSeparatorOptions(o *os.File, argv0 string) {
//	fmt.Fprintf(o, "  --rs     --irs     --ors              Record separators, e.g. 'lf' or '\\r\\n'\n");
//	fmt.Fprintf(o, "  --fs     --ifs     --ofs  --repifs    Field separators, e.g. comma\n");
//	fmt.Fprintf(o, "  --ps     --ips     --ops              Pair separators, e.g. equals sign\n");
//	fmt.Fprintf(o, "\n");
//	fmt.Fprintf(o, "  Notes about line endings:\n");
//	fmt.Fprintf(o, "  * Default line endings (--irs and --ors) are \"auto\" which means autodetect from\n");
//	fmt.Fprintf(o, "    the input file format, as long as the input file(s) have lines ending in either\n");
//	fmt.Fprintf(o, "    LF (also known as linefeed, '\\n', 0x0a, Unix-style) or CRLF (also known as\n");
//	fmt.Fprintf(o, "    carriage-return/linefeed pairs, '\\r\\n', 0x0d 0x0a, Windows style).\n");
//	fmt.Fprintf(o, "  * If both irs and ors are auto (which is the default) then LF input will lead to LF\n");
//	fmt.Fprintf(o, "    output and CRLF input will lead to CRLF output, regardless of the platform you're\n");
//	fmt.Fprintf(o, "    running on.\n");
//	fmt.Fprintf(o, "  * The line-ending autodetector triggers on the first line ending detected in the input\n");
//	fmt.Fprintf(o, "    stream. E.g. if you specify a CRLF-terminated file on the command line followed by an\n");
//	fmt.Fprintf(o, "    LF-terminated file then autodetected line endings will be CRLF.\n");
//	fmt.Fprintf(o, "  * If you use --ors {something else} with (default or explicitly specified) --irs auto\n");
//	fmt.Fprintf(o, "    then line endings are autodetected on input and set to what you specify on output.\n");
//	fmt.Fprintf(o, "  * If you use --irs {something else} with (default or explicitly specified) --ors auto\n");
//	fmt.Fprintf(o, "    then the output line endings used are LF on Unix/Linux/BSD/MacOSX, and CRLF on Windows.\n");
//	fmt.Fprintf(o, "\n");
//	fmt.Fprintf(o, "  Notes about all other separators:\n");
//	fmt.Fprintf(o, "  * IPS/OPS are only used for DKVP and XTAB formats, since only in these formats\n");
//	fmt.Fprintf(o, "    do key-value pairs appear juxtaposed.\n");
//	fmt.Fprintf(o, "  * IRS/ORS are ignored for XTAB format. Nominally IFS and OFS are newlines;\n");
//	fmt.Fprintf(o, "    XTAB records are separated by two or more consecutive IFS/OFS -- i.e.\n");
//	fmt.Fprintf(o, "    a blank line. Everything above about --irs/--ors/--rs auto becomes --ifs/--ofs/--fs\n");
//	fmt.Fprintf(o, "    auto for XTAB format. (XTAB's default IFS/OFS are \"auto\".)\n");
//	fmt.Fprintf(o, "  * OFS must be single-character for PPRINT format. This is because it is used\n");
//	fmt.Fprintf(o, "    with repetition for alignment; multi-character separators would make\n");
//	fmt.Fprintf(o, "    alignment impossible.\n");
//	fmt.Fprintf(o, "  * OPS may be multi-character for XTAB format, in which case alignment is\n");
//	fmt.Fprintf(o, "    disabled.\n");
//	fmt.Fprintf(o, "  * TSV is simply CSV using tab as field separator (\"--fs tab\").\n");
//	fmt.Fprintf(o, "  * FS/PS are ignored for markdown format; RS is used.\n");
//	fmt.Fprintf(o, "  * All FS and PS options are ignored for JSON format, since they are not relevant\n");
//	fmt.Fprintf(o, "    to the JSON format.\n");
//	fmt.Fprintf(o, "  * You can specify separators in any of the following ways, shown by example:\n");
//	fmt.Fprintf(o, "    - Type them out, quoting as necessary for shell escapes, e.g.\n");
//	fmt.Fprintf(o, "      \"--fs '|' --ips :\"\n");
//	fmt.Fprintf(o, "    - C-style escape sequences, e.g. \"--rs '\\r\\n' --fs '\\t'\".\n");
//	fmt.Fprintf(o, "    - To avoid backslashing, you can use any of the following names:\n");
//	fmt.Fprintf(o, "     ");
//	lhmss_t* pmap = get_desc_to_chars_map();
//	for (lhmsse_t* pe = pmap.phead; pe != nil; pe = pe.pnext) {
//		fmt.Fprintf(o, " %s", pe.key);
//	}
//	fmt.Fprintf(o, "\n");
//	fmt.Fprintf(o, "  * Default separators by format:\n");
//	fmt.Fprintf(o, "      %-12s %-8s %-8s %s\n", "File format", "RS", "FS", "PS");
//	lhmss_t* default_rses = get_default_rses();
//	lhmss_t* default_fses = get_default_fses();
//	lhmss_t* default_pses = get_default_pses();
//	for (lhmsse_t* pe = default_rses.phead; pe != nil; pe = pe.pnext) {
//		char* filefmt = pe.key;
//		char* rs = pe.value;
//		char* fs = lhmss_get(default_fses, filefmt);
//		char* ps = lhmss_get(default_pses, filefmt);
//		fmt.Fprintf(o, "      %-12s %-8s %-8s %s\n", filefmt, rebackslash(rs), rebackslash(fs), rebackslash(ps));
//	}
//}

//func mainUsageCsvOptions(o *os.File, argv0 string) {
//	fmt.Fprintf(o, "  --implicit-csv-header Use 1,2,3,... as field labels, rather than from line 1\n");
//	fmt.Fprintf(o, "                     of input files. Tip: combine with \"label\" to recreate\n");
//	fmt.Fprintf(o, "                     missing headers.\n");
//	fmt.Fprintf(o, "  --allow-ragged-csv-input|--ragged If a data line has fewer fields than the header line,\n");
//	fmt.Fprintf(o, "                     fill remaining keys with empty string. If a data line has more\n");
//	fmt.Fprintf(o, "                     fields than the header line, use integer field labels as in\n");
//	fmt.Fprintf(o, "                     the implicit-header case.\n");
//	fmt.Fprintf(o, "  --headerless-csv-output   Print only CSV data lines.\n");
//	fmt.Fprintf(o, "  -N                 Keystroke-saver for --implicit-csv-header --headerless-csv-output.\n");
//}

//func mainUsageDoubleQuoting(o *os.File, argv0 string) {
//	fmt.Fprintf(o, "  --quote-all        Wrap all fields in double quotes\n");
//	fmt.Fprintf(o, "  --quote-none       Do not wrap any fields in double quotes, even if they have\n");
//	fmt.Fprintf(o, "                     OFS or ORS in them\n");
//	fmt.Fprintf(o, "  --quote-minimal    Wrap fields in double quotes only if they have OFS or ORS\n");
//	fmt.Fprintf(o, "                     in them (default)\n");
//	fmt.Fprintf(o, "  --quote-numeric    Wrap fields in double quotes only if they have numbers\n");
//	fmt.Fprintf(o, "                     in them\n");
//	fmt.Fprintf(o, "  --quote-original   Wrap fields in double quotes if and only if they were\n");
//	fmt.Fprintf(o, "                     quoted on input. This isn't sticky for computed fields:\n");
//	fmt.Fprintf(o, "                     e.g. if fields a and b were quoted on input and you do\n");
//	fmt.Fprintf(o, "                     \"put '$c = $a . $b'\" then field c won't inherit a or b's\n");
//	fmt.Fprintf(o, "                     was-quoted-on-input flag.\n");
//}

//func mainUsageNumericalFormatting(o *os.File, argv0 string) {
//	fmt.Fprintf(o, "  --ofmt {format}    E.g. %%.18lf, %%.0lf. Please use sprintf-style codes for\n");
//	fmt.Fprintf(o, "                     double-precision. Applies to verbs which compute new\n");
//	fmt.Fprintf(o, "                     values, e.g. put, stats1, stats2. See also the fmtnum\n");
//	fmt.Fprintf(o, "                     function within mlr put (mlr --help-all-functions).\n");
//	fmt.Fprintf(o, "                     Defaults to %s.\n", DEFAULT_OFMT);
//}

//func mainUsageOtherOptions(o *os.File, argv0 string) {
//	fmt.Fprintf(o, "  --seed {n} with n of the form 12345678 or 0xcafefeed. For put/filter\n");
//	fmt.Fprintf(o, "                     urand()/urandint()/urand32().\n");
//	fmt.Fprintf(o, "  --nr-progress-mod {m}, with m a positive integer: print filename and record\n");
//	fmt.Fprintf(o, "                     count to os.Stderr every m input records.\n");
//	fmt.Fprintf(o, "  --from {filename}  Use this to specify an input file before the verb(s),\n");
//	fmt.Fprintf(o, "                     rather than after. May be used more than once. Example:\n");
//	fmt.Fprintf(o, "                     \"%s --from a.dat --from b.dat cat\" is the same as\n", argv0);
//	fmt.Fprintf(o, "                     \"%s cat a.dat b.dat\".\n", argv0);
//	fmt.Fprintf(o, "  -n                 Process no input files, nor standard input either. Useful\n");
//	fmt.Fprintf(o, "                     for %s put with begin/end statements only. (Same as --from\n", argv0);
//	fmt.Fprintf(o, "                     /dev/null.) Also useful in \"%s -n put -v '...'\" for\n", argv0);
//	fmt.Fprintf(o, "                     analyzing abstract syntax trees (if that's your thing).\n");
//	fmt.Fprintf(o, "  -I                 Process files in-place. For each file name on the command\n");
//	fmt.Fprintf(o, "                     line, output is written to a temp file in the same\n");
//	fmt.Fprintf(o, "                     directory, which is then renamed over the original. Each\n");
//	fmt.Fprintf(o, "                     file is processed in isolation: if the output format is\n");
//	fmt.Fprintf(o, "                     CSV, CSV headers will be present in each output file;\n");
//	fmt.Fprintf(o, "                     statistics are only over each file's own records; and so on.\n");
//}

func mainUsageThenChaining(o *os.File, argv0 string) {
	fmt.Fprintf(o, "Output of one verb may be chained as input to another using \"then\", e.g.\n")
	fmt.Fprintf(o, "  %s stats1 -a min,mean,max -f flag,u,v -g color then sort -f color\n", argv0)
}

//func mainUsageAuxents(o *os.File, argv0 string) {
//	fmt.Fprintf(o, "Miller has a few otherwise-standalone executables packaged within it.\n");
//	fmt.Fprintf(o, "They do not participate in any other parts of Miller.\n");
//	show_aux_entries(o);
//}

func mainUsageSeeAlso(o *os.File, argv0 string) {
	fmt.Fprintf(o, "For more information please see http://johnkerl.org/miller/doc and/or\n")
	fmt.Fprintf(o, "http://github.com/johnkerl/miller.")
	fmt.Fprintf(o, " This is Miller version %s.\n", VERSION_STRING)
}

//func printTypeArithmeticInfo(o *os.File, argv0 string) {
//	for (int i = -2; i < MT_DIM; i++) {
//		mv_t a = (mv_t) {.type = i, .free_flags = NO_FREE, .u.intv = 0};
//		if (i == -2)
//			fmt.Printf("%-6s |", "(+)");
//		else if (i == -1)
//			fmt.Printf("%-6s +", "------");
//		else
//			fmt.Printf("%-6s |", mt_describe_type_simple(a.type));
//
//		for (int j = 0; j < MT_DIM; j++) {
//			mv_t b = (mv_t) {.type = j, .free_flags = NO_FREE, .u.intv = 0};
//			if (i == -2) {
//				fmt.Printf(" %-6s", mt_describe_type_simple(b.type));
//			} else if (i == -1) {
//				fmt.Printf(" %-6s", "------");
//			} else {
//				mv_t c = x_xx_plus_func(&a, &b);
//				fmt.Printf(" %-6s", mt_describe_type_simple(c.type));
//			}
//		}
//
//		fmt.Fprintf(o, "\n");
//	}
//}

func usageUnrecognizedVerb(argv0 string, arg string) {
	fmt.Fprintf(os.Stderr, "%s: option \"%s\" not recognized.\n", argv0, arg)
	fmt.Fprintf(os.Stderr, "Please run \"%s --help\" for usage information.\n", argv0)
	os.Exit(1)
}

func checkArgCount(args []string, argi int, argc int, n int) {
	if (argc - argi) < n {
		fmt.Fprintf(os.Stderr, "%s: option \"%s\" missing argument(s).\n", args[0], args[argi])
		mainUsageShort()
		os.Exit(1)
	}
}

//	cli_reader_opts_init(&options.ReaderOptions);
//	cli_writer_opts_init(&options.WriterOptions);
//
//	options.mapper_argb     = 0;
//	options.filenames       = slls_alloc();
//
//	options.ofmt            = nil;
//	options.nr_progress_mod = 0LL;
//
//	options.do_in_place     = false;
//
//	options.no_input        = false;
//	options.have_rand_seed  = false;
//	options.rand_seed       = 0;
//}

// ----------------------------------------------------------------
// * If $MLRRC is set, use it and only it.
// * Otherwise try first $HOME/.mlrrc and then ./.mlrrc but let them
//   stack: e.g. $HOME/.mlrrc is lots of settings and maybe in one
//   subdir you want to override just a setting or two.

//func loadMlrrcOrDie(cli_opts_t* popts) {
//	char* env_mlrrc = getenv("MLRRC");
//	if (env_mlrrc != nil) {
//		if env_mlrrc == "__none__" {
//			return;
//		}
//		if (tryLoadMlrrc(popts, env_mlrrc)) {
//			return;
//		}
//	}
//
//	char* env_home = getenv("HOME");
//	if (env_home != nil) {
//		char* path = mlr_paste_2_strings(env_home, "/.mlrrc");
//		(void)tryLoadMlrrc(popts, path);
//		free(path);
//	}
//
//	(void)tryLoadMlrrc(popts, "./.mlrrc");
//}

//static int tryLoadMlrrc(cli_opts_t* popts, char* path) {
//	FILE* fp = fopen(path, "r");
//	if (fp == nil) {
//		return false;
//	}
//
//	char* line = nil;
//	size_t linecap = 0;
//	int rc;
//	int lineno = 0;
//
//	while ((rc = getline(&line, &linecap, fp)) != -1) {
//		lineno++;
//		char* line_to_destroy = strdup(line);
//		if (!handle_mlrrc_line_1(popts, line_to_destroy)) {
//			fmt.Fprintf(os.Stderr, "Parse error at file \"%s\" line %d: %s\n",
//				path, lineno, line);
//			os.Exit(1);
//		}
//		free(line_to_destroy);
//	}
//
//	fclose(fp);
//	if (line != nil) {
//		free(line);
//	}
//
//	return true;
//}

// Chomps trailing CR, LF, or CR/LF; comment-strips; left-right trims.

//static int handle_mlrrc_line_1(cli_opts_t* popts, char* line) {
//	// chomp
//	size_t len = len(line);
//	if (len >= 2 && line[len-2] == '\r' && line[len-1] == '\n') {
//		line[len-2] = 0;
//	} else if (len >= 1 && (line[len-1] == '\r' || line[len-1] == '\n')) {
//		line[len-1] = 0;
//	}
//
//	// comment-strip
//	char* pbang = strstr(line, "#");
//	if (pbang != nil) {
//		*pbang = 0;
//	}
//
//	// Left-trim
//	char* start = line;
//	while (*start == ' ' || *start == '\t') {
//		start++;
//	}
//
//	// Right-trim
//	len = len(start);
//	char* end = &start[len-1];
//	while (end > start && (*end == ' ' || *end == '\t')) {
//		*end = 0;
//		end--;
//	}
//	if (end < start) { // line was whitespace-only
//		return true;
//	} else {
//		return handle_mlrrc_line_2(popts, start);
//	}
//}

// Prepends initial "--" if it's not already there
//static int handle_mlrrc_line_2(cli_opts_t* popts, char* line) {
//	size_t len = len(line);
//
//	char* dashed_line = nil;
//	if (len >= 2 && line[0] != '-' && line[1] != '-') {
//		dashed_line = mlr_paste_2_strings("--", line);
//	} else {
//		dashed_line = strdup(line);
//	}
//
//	int rc = handle_mlrrc_line_3(popts, dashed_line);
//
//	// Do not free these. The command-line parsers can retain pointers into args strings (rather
//	// than copying), resulting in freed-memory reads later in the data-processing verbs.
//	//
//	// It would be possible to be diligent about making sure all current command-line-parsing
//	// callsites copy strings rather than pointing to them -- but it would be easy to miss some, and
//	// also any future codemods might make the same mistake as well.
//	//
//	// It's safer (and no big leak) to simply leave these parsed mlrrc lines unfreed.
//	//
//	// free(dashed_line);
//	return rc;
//}

// Splits line into args array
//static int handle_mlrrc_line_3(cli_opts_t* popts, char* line) {
//	char* args[3];
//	int argc = 0;
//	char* split = strpbrk(line, " \t");
//	if (split == nil) {
//		args[0] = line;
//		args[1] = nil;
//		argc = 1;
//	} else {
//		*split = 0;
//		char* p = split + 1;
//		while (*p == ' ' || *p == '\t') {
//			p++;
//		}
//		args[0] = line;
//		args[1] = p;
//		args[2] = nil;
//		argc = 2;
//	}
//	return handle_mlrrc_line_4(popts, args, argc);
//}

//static int handle_mlrrc_line_4(cli_opts_t* popts, char** args, int argc) {
//	int argi = 0;
//	if (handleReaderOptions(args, argc, &argi, &options.ReaderOptions)) {
//		// handled
//	} else if (handleWriterOptions(args, argc, &argi, &options.WriterOptions)) {
//		// handled
//	} else if (handleReaderWriterOptions(args, argc, &argi, &options.ReaderOptions, &options.WriterOptions)) {
//		// handled
//	} else if (handleMiscOptions(args, argc, &argi, popts)) {
//		// handled
//	} else {
//		// unhandled
//		return false;
//	}
//
//	return true;
//}

// ----------------------------------------------------------------
//void cli_reader_opts_init(clitypes.TReaderOptions* readerOptions) {
//	readerOptions.InputFileFormat                      = nil;
//	readerOptions.IRS                            = nil;
//	readerOptions.IFS                            = nil;
//	readerOptions.ips                            = nil;
//	readerOptions.input_json_flatten_separator   = nil;
//	readerOptions.json_array_ingest              = JSON_ARRAY_INGEST_UNSPECIFIED;
//
//	readerOptions.allow_repeat_ifs               = NEITHER_TRUE_NOR_FALSE;
//	readerOptions.allow_repeat_ips               = NEITHER_TRUE_NOR_FALSE;
//	readerOptions.use_implicit_csv_header        = NEITHER_TRUE_NOR_FALSE;
//	readerOptions.allow_ragged_csv_input         = NEITHER_TRUE_NOR_FALSE;
//
//	readerOptions.prepipe                        = nil;
//	readerOptions.comment_handling               = COMMENTS_ARE_DATA;
//	readerOptions.comment_string                 = nil;
//
//	readerOptions.generator_opts.field_name     = "i";
//	readerOptions.generator_opts.start          = 0LL;
//	readerOptions.generator_opts.stop           = 100LL;
//	readerOptions.generator_opts.step           = 1LL;
//}

//void cli_writer_opts_init(clitypes.TWriterOptions* writerOptions) {
//	writerOptions.OutputFileFormat                      = nil;
//	writerOptions.ORS                            = nil;
//	writerOptions.OFS                            = nil;
//	writerOptions.ops                            = nil;
//
//	writerOptions.headerless_csv_output          = NEITHER_TRUE_NOR_FALSE;
//	writerOptions.right_justify_xtab_value       = NEITHER_TRUE_NOR_FALSE;
//	writerOptions.right_align_pprint             = NEITHER_TRUE_NOR_FALSE;
//	writerOptions.pprint_barred                  = NEITHER_TRUE_NOR_FALSE;
//	writerOptions.stack_json_output_vertically   = NEITHER_TRUE_NOR_FALSE;
//	writerOptions.wrap_json_output_in_outer_list = NEITHER_TRUE_NOR_FALSE;
//	writerOptions.json_quote_int_keys            = NEITHER_TRUE_NOR_FALSE;
//	writerOptions.json_quote_non_string_values   = NEITHER_TRUE_NOR_FALSE;
//
//	writerOptions.output_json_flatten_separator  = nil;
//	writerOptions.oosvar_flatten_separator       = nil;
//
//	writerOptions.oquoting                       = QUOTE_UNSPECIFIED;
//}

//void cli_apply_defaults(cli_opts_t* popts) {
//
//	cli_apply_reader_defaults(&options.ReaderOptions);
//
//	cli_apply_writer_defaults(&options.WriterOptions);
//
//	if (options.ofmt == nil)
//		options.ofmt = DEFAULT_OFMT;
//}

//void cli_apply_reader_defaults(clitypes.TReaderOptions* readerOptions) {
//	if (readerOptions.InputFileFormat == nil)
//		readerOptions.InputFileFormat = "dkvp";
//
//	if (readerOptions.json_array_ingest == JSON_ARRAY_INGEST_UNSPECIFIED)
//		readerOptions.json_array_ingest = JSON_ARRAY_INGEST_AS_MAP;
//
//	if (readerOptions.use_implicit_csv_header == NEITHER_TRUE_NOR_FALSE)
//		readerOptions.use_implicit_csv_header = false;
//
//	if (readerOptions.allow_ragged_csv_input == NEITHER_TRUE_NOR_FALSE)
//		readerOptions.allow_ragged_csv_input = false;
//
//	if (readerOptions.input_json_flatten_separator == nil)
//		readerOptions.input_json_flatten_separator = DEFAULT_JSON_FLATTEN_SEPARATOR;
//}

//void cli_apply_writer_defaults(clitypes.TWriterOptions* writerOptions) {
//	if (writerOptions.OutputFileFormat == nil)
//		writerOptions.OutputFileFormat = "dkvp";
//
//	if (writerOptions.headerless_csv_output == NEITHER_TRUE_NOR_FALSE)
//		writerOptions.headerless_csv_output = false;
//
//	if (writerOptions.right_justify_xtab_value == NEITHER_TRUE_NOR_FALSE)
//		writerOptions.right_justify_xtab_value = false;
//
//	if (writerOptions.right_align_pprint == NEITHER_TRUE_NOR_FALSE)
//		writerOptions.right_align_pprint = false;
//
//	if (writerOptions.pprint_barred == NEITHER_TRUE_NOR_FALSE)
//		writerOptions.pprint_barred = false;
//
//	if (writerOptions.stack_json_output_vertically == NEITHER_TRUE_NOR_FALSE)
//		writerOptions.stack_json_output_vertically = false;
//
//	if (writerOptions.wrap_json_output_in_outer_list == NEITHER_TRUE_NOR_FALSE)
//		writerOptions.wrap_json_output_in_outer_list = false;
//
//	if (writerOptions.json_quote_int_keys == NEITHER_TRUE_NOR_FALSE)
//		writerOptions.json_quote_int_keys = true;
//
//	if (writerOptions.json_quote_non_string_values == NEITHER_TRUE_NOR_FALSE)
//		writerOptions.json_quote_non_string_values = false;
//
//	if (writerOptions.output_json_flatten_separator == nil)
//		writerOptions.output_json_flatten_separator = DEFAULT_JSON_FLATTEN_SEPARATOR;
//
//	if (writerOptions.oosvar_flatten_separator == nil)
//		writerOptions.oosvar_flatten_separator = DEFAULT_OOSVAR_FLATTEN_SEPARATOR;
//
//	if (writerOptions.oquoting == QUOTE_UNSPECIFIED)
//		writerOptions.oquoting = DEFAULT_OQUOTING;
//}

// ----------------------------------------------------------------
// For mapper join which has its own input-format overrides.
//
// Mainly this just takes the main-opts flag whenever the join-opts flag was not
// specified by the user. But it's a bit more complex when main and join input
// formats are different. Example: main input format is CSV, for which IPS is
// "(N/A)", and join input format is DKVP. Then we should not use "(N/A)"
// for DKVP IPS. However if main input format were DKVP with IPS set to ":",
// then we should take that.
//
// The logic is:
//
// * If the join input format was unspecified, take all unspecified values from
//   main opts.
//
// * If the join input format was specified and is the same as main input
//   format, take unspecified values from main opts.
//
// * If the join input format was specified and is not the same as main input
//   format, take unspecified values from defaults for the join input format.

//void cli_merge_reader_opts(clitypes.TReaderOptions* pfunc_opts, TReaderOptions* pmain_opts) {
//
//	if (pfunc_opts.InputFileFormat == nil) {
//		pfunc_opts.InputFileFormat = pmain_opts.InputFileFormat;
//	}
//
//	if pfunc_opts.InputFileFormat == pmain_opts.InputFileFormat {
//
//		if (pfunc_opts.IRS == nil)
//			pfunc_opts.IRS = pmain_opts.IRS;
//		if (pfunc_opts.IFS == nil)
//			pfunc_opts.IFS = pmain_opts.IFS;
//		if (pfunc_opts.ips == nil)
//			pfunc_opts.ips = pmain_opts.ips;
//		if (pfunc_opts.allow_repeat_ifs  == NEITHER_TRUE_NOR_FALSE)
//			pfunc_opts.allow_repeat_ifs = pmain_opts.allow_repeat_ifs;
//		if (pfunc_opts.allow_repeat_ips  == NEITHER_TRUE_NOR_FALSE)
//			pfunc_opts.allow_repeat_ips = pmain_opts.allow_repeat_ips;
//
//	} else {
//
//		if (pfunc_opts.IRS == nil)
//			pfunc_opts.IRS = lhmss_get_or_die(get_default_rses(), pfunc_opts.InputFileFormat);
//		if (pfunc_opts.IFS == nil)
//			pfunc_opts.IFS = lhmss_get_or_die(get_default_fses(), pfunc_opts.InputFileFormat);
//		if (pfunc_opts.ips == nil)
//			pfunc_opts.ips = lhmss_get_or_die(get_default_pses(), pfunc_opts.InputFileFormat);
//		if (pfunc_opts.allow_repeat_ifs  == NEITHER_TRUE_NOR_FALSE)
//			pfunc_opts.allow_repeat_ifs = lhmsll_get_or_die(get_default_repeat_ifses(), pfunc_opts.InputFileFormat);
//		if (pfunc_opts.allow_repeat_ips  == NEITHER_TRUE_NOR_FALSE)
//			pfunc_opts.allow_repeat_ips = lhmsll_get_or_die(get_default_repeat_ipses(), pfunc_opts.InputFileFormat);
//
//	}
//
//	if (pfunc_opts.json_array_ingest == JSON_ARRAY_INGEST_UNSPECIFIED)
//		pfunc_opts.json_array_ingest = pmain_opts.json_array_ingest;
//
//	if (pfunc_opts.use_implicit_csv_header == NEITHER_TRUE_NOR_FALSE)
//		pfunc_opts.use_implicit_csv_header = pmain_opts.use_implicit_csv_header;
//
//	if (pfunc_opts.allow_ragged_csv_input == NEITHER_TRUE_NOR_FALSE)
//		pfunc_opts.allow_ragged_csv_input = pmain_opts.allow_ragged_csv_input;
//
//	if (pfunc_opts.input_json_flatten_separator == nil)
//		pfunc_opts.input_json_flatten_separator = pmain_opts.input_json_flatten_separator;
//}

// Similar to cli_merge_reader_opts but for mapper tee & mapper put which have their
// own output-format overrides.

//void cli_merge_writer_opts(clitypes.TWriterOptions* pfunc_opts, TWriterOptions* pmain_opts) {
//
//	if (pfunc_opts.OutputFileFormat == nil) {
//		pfunc_opts.OutputFileFormat = pmain_opts.OutputFileFormat;
//	}
//
//	if pfunc_opts.OutputFileFormat == pmain_opts.OutputFileFormat {
//		if (pfunc_opts.ORS == nil)
//			pfunc_opts.ORS = pmain_opts.ORS;
//		if (pfunc_opts.OFS == nil)
//			pfunc_opts.OFS = pmain_opts.OFS;
//		if (pfunc_opts.ops == nil)
//			pfunc_opts.ops = pmain_opts.ops;
//	} else {
//		if (pfunc_opts.ORS == nil)
//			pfunc_opts.ORS = lhmss_get_or_die(get_default_rses(), pfunc_opts.OutputFileFormat);
//		if (pfunc_opts.OFS == nil)
//			pfunc_opts.OFS = lhmss_get_or_die(get_default_fses(), pfunc_opts.OutputFileFormat);
//		if (pfunc_opts.ops == nil)
//			pfunc_opts.ops = lhmss_get_or_die(get_default_pses(), pfunc_opts.OutputFileFormat);
//	}
//
//	if (pfunc_opts.headerless_csv_output == NEITHER_TRUE_NOR_FALSE)
//		pfunc_opts.headerless_csv_output = pmain_opts.headerless_csv_output;
//
//	if (pfunc_opts.right_justify_xtab_value == NEITHER_TRUE_NOR_FALSE)
//		pfunc_opts.right_justify_xtab_value = pmain_opts.right_justify_xtab_value;
//
//	if (pfunc_opts.right_align_pprint == NEITHER_TRUE_NOR_FALSE)
//		pfunc_opts.right_align_pprint = pmain_opts.right_align_pprint;
//
//	if (pfunc_opts.pprint_barred == NEITHER_TRUE_NOR_FALSE)
//		pfunc_opts.pprint_barred = pmain_opts.pprint_barred;
//
//	if (pfunc_opts.stack_json_output_vertically == NEITHER_TRUE_NOR_FALSE)
//		pfunc_opts.stack_json_output_vertically = pmain_opts.stack_json_output_vertically;
//
//	if (pfunc_opts.wrap_json_output_in_outer_list == NEITHER_TRUE_NOR_FALSE)
//		pfunc_opts.wrap_json_output_in_outer_list = pmain_opts.wrap_json_output_in_outer_list;
//
//	if (pfunc_opts.json_quote_int_keys == NEITHER_TRUE_NOR_FALSE)
//		pfunc_opts.json_quote_int_keys = pmain_opts.json_quote_int_keys;
//
//	if (pfunc_opts.json_quote_non_string_values == NEITHER_TRUE_NOR_FALSE)
//		pfunc_opts.json_quote_non_string_values = pmain_opts.json_quote_non_string_values;
//
//	if (pfunc_opts.output_json_flatten_separator == nil)
//		pfunc_opts.output_json_flatten_separator = pmain_opts.output_json_flatten_separator;
//
//	if (pfunc_opts.oosvar_flatten_separator == nil)
//		pfunc_opts.oosvar_flatten_separator = pmain_opts.oosvar_flatten_separator;
//
//	if (pfunc_opts.oquoting == QUOTE_UNSPECIFIED)
//		pfunc_opts.oquoting = pmain_opts.oquoting;
//}

// ----------------------------------------------------------------
func handleTerminalUsage(args []string, argc int, argi int) bool {
	if args[argi] == "--version" {
		fmt.Printf("Miller %s\n", VERSION_STRING)
		return true
	} else if args[argi] == "-h" {
		mainUsageLong(os.Stdout, os.Args[0])
		return true
	} else if args[argi] == "--help" {
		mainUsageLong(os.Stdout, os.Args[0])
		return true
		//	} else if args[argi] == "--print-type-arithmetic-info" {
		//		printTypeArithmeticInfo(os.Stdout, os.Args[0]);
		//		return true;
		//
	} else if args[argi] == "--help-all-verbs" || args[argi] == "--usage-all-verbs" {
		usageAllVerbs(os.Args[0])
	} else if args[argi] == "--list-all-verbs" || args[argi] == "-l" {
		listAllVerbs(os.Stdout, " ")
		return true
	} else if args[argi] == "--list-all-verbs-raw" || args[argi] == "-L" {
		listAllVerbsRaw(os.Stdout)
		return true

		//	} else if args[argi] == "--list-all-functions-raw" || args[argi] == "-F" {
		//		fmgr_t* pfmgr = fmgr_alloc();
		//		fmgr_list_all_functions_raw(pfmgr, os.Stdout);
		//		fmgr_free(pfmgr, nil);
		//		return true;
		//	} else if args[argi] == "--list-all-functions-as-table" {
		//		fmgr_t* pfmgr = fmgr_alloc();
		//		fmgr_list_all_functions_as_table(pfmgr, os.Stdout);
		//		fmgr_free(pfmgr, nil);
		//		return true;
		//	} else if args[argi] == "--help-all-functions" || args[argi] == "-f" {
		//		fmgr_t* pfmgr = fmgr_alloc();
		//		fmgr_function_usage(pfmgr, os.Stdout, nil);
		//		fmgr_free(pfmgr, nil);
		//		return true;
		//	} else if args[argi] == "--help-function" || args[argi] == "--hf" {
		//		checkArgCount(args, argi, argc, 2);
		//		fmgr_t* pfmgr = fmgr_alloc();
		//		fmgr_function_usage(pfmgr, os.Stdout, args[argi+1]);
		//		fmgr_free(pfmgr, nil);
		//		return true;
		//
		//	} else if args[argi] == "--list-all-keywords-raw" || args[argi] == "-K" {
		//		mlr_dsl_list_all_keywords_raw(os.Stdout);
		//		return true;
		//	} else if args[argi] == "--help-all-keywords" || args[argi] == "-k" {
		//		mlr_dsl_keyword_usage(os.Stdout, nil);
		//		return true;
		//	} else if args[argi] == "--help-keyword" || args[argi] == "--hk" {
		//		checkArgCount(args, argi, argc, 2);
		//		mlr_dsl_keyword_usage(os.Stdout, args[argi+1]);
		//		return true;
		//
		//	// main-usage subsections, individually accessible for the benefit of
		//	// the manpage-autogenerator
	} else if args[argi] == "--usage-synopsis" {
		mainUsageSynopsis(os.Stdout, os.Args[0])
		return true
	} else if args[argi] == "--usage-examples" {
		mainUsageExamples(os.Stdout, os.Args[0], "")
		return true
	} else if args[argi] == "--usage-list-all-verbs" {
		listAllVerbs(os.Stdout, "")
		return true
	} else if args[argi] == "--usage-help-options" {
		mainUsageHelpOptions(os.Stdout, os.Args[0])
		return true
		//	} else if args[argi] == "--usage-mlrrc" {
		//		mainUsageMlrrc(os.Stdout, os.Args[0]);
		//		return true;
		//	} else if args[argi] == "--usage-functions" {
		//		mainUsageFunctions(os.Stdout, os.Args[0], "");
		//		return true;
		//	} else if args[argi] == "--usage-data-format-examples" {
		mainUsageDataFormatExamples(os.Stdout, os.Args[0])
		//		return true;
	} else if args[argi] == "--usage-data-format-options" {
		mainUsageDataFormatOptions(os.Stdout, os.Args[0])
		return true
		//	} else if args[argi] == "--usage-comments-in-data" {
		//		mainUsageCommentsInData(os.Stdout, os.Args[0]);
		//		return true;
	} else if args[argi] == "--usage-format-conversion-keystroke-saver-options" {
		mainUsageFormatConversionKeystrokeSaverOptions(os.Stdout, os.Args[0])
		return true
		//	} else if args[argi] == "--usage-compressed-data-options" {
		//		mainUsageCompressedDataOptions(os.Stdout, os.Args[0]);
		//		return true;
		//	} else if args[argi] == "--usage-separator-options" {
		//		mainUsageSeparatorOptions(os.Stdout, os.Args[0]);
		//		return true;
		//	} else if args[argi] == "--usage-csv-options" {
		//		mainUsageCsvOptions(os.Stdout, os.Args[0]);
		//		return true;
		//	} else if args[argi] == "--usage-double-quoting" {
		//		mainUsageDoubleQuoting(os.Stdout, os.Args[0]);
		//		return true;
		//	} else if args[argi] == "--usage-numerical-formatting" {
		//		mainUsageNumericalFormatting(os.Stdout, os.Args[0]);
		//		return true;
		//	} else if args[argi] == "--usage-other-options" {
		//		mainUsageOtherOptions(os.Stdout, os.Args[0]);
		//		return true;
	} else if args[argi] == "--usage-then-chaining" {
		mainUsageThenChaining(os.Stdout, os.Args[0])
		return true
		//	} else if args[argi] == "--usage-auxents" {
		//		mainUsageAuxents(os.Stdout, os.Args[0]);
		//		return true;
	} else if args[argi] == "--usage-see-also" {
		mainUsageSeeAlso(os.Stdout, os.Args[0])
		return true
	}
	return false
}

// Returns true if the current flag was handled.
func handleReaderOptions(args []string, argc int, pargi *int, readerOptions *clitypes.TReaderOptions) bool {
	argi := *pargi
	oargi := argi

	if args[argi] == "--irs" {
		checkArgCount(args, argi, argc, 2)
		readerOptions.IRS = SeparatorFromArg(args[argi+1])
		argi += 2

	} else if args[argi] == "--ifs" {
		checkArgCount(args, argi, argc, 2)
		readerOptions.IFS = SeparatorFromArg(args[argi+1])
		argi += 2

	} else if args[argi] == "--ips" {
		checkArgCount(args, argi, argc, 2)
		readerOptions.IPS = SeparatorFromArg(args[argi+1])
		argi += 2

		//	} else if args[argi] == "--repifs" {
		//		readerOptions.allow_repeat_ifs = true;
		//		argi += 1;
		//
		//	} else if args[argi] == "--json-fatal-arrays-on-input" {
		//		readerOptions.json_array_ingest = JSON_ARRAY_INGEST_FATAL;
		//		argi += 1;
		//	} else if args[argi] == "--json-skip-arrays-on-input" {
		//		readerOptions.json_array_ingest = JSON_ARRAY_INGEST_SKIP;
		//		argi += 1;
		//	} else if args[argi] == "--json-map-arrays-on-input" {
		//		readerOptions.json_array_ingest = JSON_ARRAY_INGEST_AS_MAP;
		//		argi += 1;
		//
		//	} else if args[argi] == "--implicit-csv-header" {
		//		readerOptions.use_implicit_csv_header = true;
		//		argi += 1;
		//
		//	} else if args[argi] == "--no-implicit-csv-header" {
		//		readerOptions.use_implicit_csv_header = false;
		//		argi += 1;
		//
		//	} else if args[argi] == "--allow-ragged-csv-input" || args[argi] == "--ragged" {
		//		readerOptions.allow_ragged_csv_input = true;
		//		argi += 1;

		//	} else if args[argi] == "-i" {
		//		checkArgCount(args, argi, argc, 2);
		//		if (!lhmss_has_key(get_default_rses(), args[argi+1])) {
		//			fmt.Fprintf(os.Stderr, "%s: unrecognized input format \"%s\".\n",
		//				os.Args[0], args[argi+1]);
		//			os.Exit(1);
		//		}
		//		readerOptions.InputFileFormat = args[argi+1];
		//		argi += 2;
		//
		//	} else if args[argi] == "--igen" {
		//		readerOptions.InputFileFormat = "gen";
		//		argi += 1;
		//	} else if args[argi] == "--gen-start" {
		//		readerOptions.InputFileFormat = "gen";
		//		checkArgCount(args, argi, argc, 2);
		//		if (sscanf(args[argi+1], "%lld", &readerOptions.generator_opts.start) != 1) {
		//			fmt.Fprintf(os.Stderr, "%s: could not scan \"%s\".\n",
		//				os.Args[0], args[argi+1]);
		//		}
		//		argi += 2;
		//	} else if args[argi] == "--gen-stop" {
		//		readerOptions.InputFileFormat = "gen";
		//		checkArgCount(args, argi, argc, 2);
		//		if (sscanf(args[argi+1], "%lld", &readerOptions.generator_opts.stop) != 1) {
		//			fmt.Fprintf(os.Stderr, "%s: could not scan \"%s\".\n",
		//				os.Args[0], args[argi+1]);
		//		}
		//		argi += 2;
		//	} else if args[argi] == "--gen-step" {
		//		readerOptions.InputFileFormat = "gen";
		//		checkArgCount(args, argi, argc, 2);
		//		if (sscanf(args[argi+1], "%lld", &readerOptions.generator_opts.step) != 1) {
		//			fmt.Fprintf(os.Stderr, "%s: could not scan \"%s\".\n",
		//				os.Args[0], args[argi+1]);
		//		}
		//		argi += 2;

	} else if args[argi] == "--icsv" {
		readerOptions.InputFileFormat = "csv"
		argi += 1

		//	} else if args[argi] == "--icsvlite" {
		//		readerOptions.InputFileFormat = "csvlite";
		//		argi += 1;
		//
		//	} else if args[argi] == "--itsv" {
		//		readerOptions.InputFileFormat = "csv";
		//		readerOptions.IFS = "\t";
		//		argi += 1;
		//
		//	} else if args[argi] == "--itsvlite" {
		//		readerOptions.InputFileFormat = "csvlite";
		//		readerOptions.IFS = "\t";
		//		argi += 1;
		//
		//	} else if args[argi] == "--iasv" {
		//		readerOptions.InputFileFormat = "csv";
		//		readerOptions.IFS = ASV_FS;
		//		readerOptions.IRS = ASV_RS;
		//		argi += 1;
		//
		//	} else if args[argi] == "--iasvlite" {
		//		readerOptions.InputFileFormat = "csvlite";
		//		readerOptions.IFS = ASV_FS;
		//		readerOptions.IRS = ASV_RS;
		//		argi += 1;
		//
		//	} else if args[argi] == "--iusv" {
		//		readerOptions.InputFileFormat = "csv";
		//		readerOptions.IFS = USV_FS;
		//		readerOptions.IRS = USV_RS;
		//		argi += 1;
		//
		//	} else if args[argi] == "--iusvlite" {
		//		readerOptions.InputFileFormat = "csvlite";
		//		readerOptions.IFS = USV_FS;
		//		readerOptions.IRS = USV_RS;
		//		argi += 1;
		//
	} else if args[argi] == "--idkvp" {
		readerOptions.InputFileFormat = "dkvp"
		argi += 1

	} else if args[argi] == "--ijson" {
		readerOptions.InputFileFormat = "json"
		argi += 1

	} else if args[argi] == "--inidx" {
		readerOptions.InputFileFormat = "nidx"
		argi += 1

		//	} else if args[argi] == "--ixtab" {
		//		readerOptions.InputFileFormat = "xtab";
		//		argi += 1;
		//
		//	} else if args[argi] == "--ipprint" {
		//		readerOptions.InputFileFormat        = "csvlite";
		//		readerOptions.IFS              = " ";
		//		readerOptions.allow_repeat_ifs = true;
		//		argi += 1;
		//
		//	} else if args[argi] == "--mmap" {
		//		// No-op as of 5.6.3 (mmap is being abandoned) but don't break
		//		// the command-line user experience.
		//		argi += 1;
		//
		//	} else if args[argi] == "--no-mmap" {
		//		// No-op as of 5.6.3 (mmap is being abandoned) but don't break
		//		// the command-line user experience.
		//		argi += 1;
		//
		//	} else if args[argi] == "--prepipe" {
		//		checkArgCount(args, argi, argc, 2);
		//		readerOptions.prepipe = args[argi+1];
		//		argi += 2;
		//
		//	} else if args[argi] == "--skip-comments" {
		//		readerOptions.comment_string = DEFAULT_COMMENT_STRING;
		//		readerOptions.comment_handling = SKIP_COMMENTS;
		//		argi += 1;
		//
		//	} else if args[argi] == "--skip-comments-with" {
		//		checkArgCount(args, argi, argc, 2);
		//		readerOptions.comment_string = args[argi+1];
		//		readerOptions.comment_handling = SKIP_COMMENTS;
		//		argi += 2;
		//
		//	} else if args[argi] == "--pass-comments" {
		//		readerOptions.comment_string = DEFAULT_COMMENT_STRING;
		//		readerOptions.comment_handling = PASS_COMMENTS;
		//		argi += 1;
		//
		//	} else if args[argi] == "--pass-comments-with" {
		//		checkArgCount(args, argi, argc, 2);
		//		readerOptions.comment_string = args[argi+1];
		//		readerOptions.comment_handling = PASS_COMMENTS;
		//		argi += 2;
		//
	}
	*pargi = argi
	return argi != oargi
}

// Returns true if the current flag was handled.
func handleWriterOptions(args []string, argc int, pargi *int, writerOptions *clitypes.TWriterOptions) bool {
	argi := *pargi
	oargi := argi

	if args[argi] == "--ors" {
		checkArgCount(args, argi, argc, 2)
		writerOptions.ORS = SeparatorFromArg(args[argi+1])
		argi += 2

	} else if args[argi] == "--ofs" {
		checkArgCount(args, argi, argc, 2)
		writerOptions.OFS = SeparatorFromArg(args[argi+1])
		argi += 2

		//	} else if args[argi] == "--headerless-csv-output" {
		//		writerOptions.headerless_csv_output = true;
		//		argi += 1;
		//
	} else if args[argi] == "--ops" {
		checkArgCount(args, argi, argc, 2)
		writerOptions.OPS = SeparatorFromArg(args[argi+1])
		argi += 2

		//	} else if args[argi] == "--xvright" {
		//		writerOptions.right_justify_xtab_value = true;
		//		argi += 1;
		//
		//	} else if args[argi] == "--jvstack" {
		//		writerOptions.stack_json_output_vertically = true;
		//		argi += 1;
		//
		//	} else if args[argi] == "--jlistwrap" {
		//		writerOptions.wrap_json_output_in_outer_list = true;
		//		argi += 1;
		//
		//	} else if args[argi] == "--jknquoteint" {
		//		writerOptions.json_quote_int_keys = false;
		//		argi += 1;
		//	} else if args[argi] == "--jquoteall" {
		//		writerOptions.json_quote_non_string_values = true;
		//		argi += 1;
		//	} else if args[argi] == "--jvquoteall" {
		//		writerOptions.json_quote_non_string_values = true;
		//		argi += 1;
		//
		//	} else if args[argi] == "--vflatsep" {
		//		checkArgCount(args, argi, argc, 2);
		//		writerOptions.oosvar_flatten_separator = SeparatorFromArg(args[argi+1]);
		//		argi += 2;
		//
		//	} else if args[argi] == "-o" {
		//		checkArgCount(args, argi, argc, 2);
		//		if (!lhmss_has_key(get_default_rses(), args[argi+1])) {
		//			fmt.Fprintf(os.Stderr, "%s: unrecognized output format \"%s\".\n",
		//				os.Args[0], args[argi+1]);
		//			os.Exit(1);
		//		}
		//		writerOptions.OutputFileFormat = args[argi+1];
		//		argi += 2;
		//
	} else if args[argi] == "--ocsv" {
		writerOptions.OutputFileFormat = "csv"
		argi += 1

		//	} else if args[argi] == "--ocsvlite" {
		//		writerOptions.OutputFileFormat = "csvlite";
		//		argi += 1;
		//
		//	} else if args[argi] == "--otsv" {
		//		writerOptions.OutputFileFormat = "csv";
		//		writerOptions.OFS = "\t";
		//		argi += 1;
		//
		//	} else if args[argi] == "--otsvlite" {
		//		writerOptions.OutputFileFormat = "csvlite";
		//		writerOptions.OFS = "\t";
		//		argi += 1;
		//
		//	} else if args[argi] == "--oasv" {
		//		writerOptions.OutputFileFormat = "csv";
		//		writerOptions.OFS = ASV_FS;
		//		writerOptions.ORS = ASV_RS;
		//		argi += 1;
		//
		//	} else if args[argi] == "--oasvlite" {
		//		writerOptions.OutputFileFormat = "csvlite";
		//		writerOptions.OFS = ASV_FS;
		//		writerOptions.ORS = ASV_RS;
		//		argi += 1;
		//
		//	} else if args[argi] == "--ousv" {
		//		writerOptions.OutputFileFormat = "csv";
		//		writerOptions.OFS = USV_FS;
		//		writerOptions.ORS = USV_RS;
		//		argi += 1;
		//
		//	} else if args[argi] == "--ousvlite" {
		//		writerOptions.OutputFileFormat = "csvlite";
		//		writerOptions.OFS = USV_FS;
		//		writerOptions.ORS = USV_RS;
		//		argi += 1;
		//
		//	} else if args[argi] == "--omd" {
		//		writerOptions.OutputFileFormat = "markdown";
		//		argi += 1;
		//
	} else if args[argi] == "--odkvp" {
		writerOptions.OutputFileFormat = "dkvp"
		argi += 1

	} else if args[argi] == "--ojson" {
		writerOptions.OutputFileFormat = "json"
		argi += 1
		//	} else if args[argi] == "--ojsonx" {
		//		writerOptions.OutputFileFormat = "json";
		//		writerOptions.stack_json_output_vertically = true;
		//		argi += 1;

	} else if args[argi] == "--onidx" {
		writerOptions.OutputFileFormat = "nidx"
		argi += 1

	} else if args[argi] == "--oxtab" {
		writerOptions.OutputFileFormat = "xtab"
		argi += 1

	} else if args[argi] == "--opprint" {
		writerOptions.OutputFileFormat = "pprint"
		argi += 1

		//	} else if args[argi] == "--right" {
		//		writerOptions.right_align_pprint = true;
		//		argi += 1;
		//
		//	} else if args[argi] == "--barred" {
		//		writerOptions.pprint_barred = true;
		//		argi += 1;
		//
		//	} else if args[argi] == "--quote-all" {
		//		writerOptions.oquoting = QUOTE_ALL;
		//		argi += 1;
		//
		//	} else if args[argi] == "--quote-none" {
		//		writerOptions.oquoting = QUOTE_NONE;
		//		argi += 1;
		//
		//	} else if args[argi] == "--quote-minimal" {
		//		writerOptions.oquoting = QUOTE_MINIMAL;
		//		argi += 1;
		//
		//	} else if args[argi] == "--quote-numeric" {
		//		writerOptions.oquoting = QUOTE_NUMERIC;
		//		argi += 1;
		//
		//	} else if args[argi] == "--quote-original" {
		//		writerOptions.oquoting = QUOTE_ORIGINAL;
		//		argi += 1;
		//
	}
	*pargi = argi
	return argi != oargi
}

// Returns true if the current flag was handled.
func handleReaderWriterOptions(
	args []string,
	argc int,
	pargi *int,
	readerOptions *clitypes.TReaderOptions,
	writerOptions *clitypes.TWriterOptions,
) bool {

	argi := *pargi
	oargi := argi

	if args[argi] == "--rs" {
		checkArgCount(args, argi, argc, 2)
		readerOptions.IRS = SeparatorFromArg(args[argi+1])
		writerOptions.ORS = SeparatorFromArg(args[argi+1])
		argi += 2

	} else if args[argi] == "--fs" {
		checkArgCount(args, argi, argc, 2)
		readerOptions.IFS = SeparatorFromArg(args[argi+1])
		writerOptions.OFS = SeparatorFromArg(args[argi+1])
		argi += 2

		//	} else if args[argi] == "-p" {
		//		readerOptions.InputFileFormat = "nidx";
		//		writerOptions.OutputFileFormat = "nidx";
		//		readerOptions.IFS = " ";
		//		writerOptions.OFS = " ";
		//		readerOptions.allow_repeat_ifs = true;
		//		argi += 1;
		//
	} else if args[argi] == "--ps" {
		checkArgCount(args, argi, argc, 2)
		readerOptions.IPS = SeparatorFromArg(args[argi+1])
		writerOptions.OPS = SeparatorFromArg(args[argi+1])
		argi += 2

		//	} else if args[argi] == "--jflatsep" {
		//		checkArgCount(args, argi, argc, 2);
		//		readerOptions.input_json_flatten_separator  = SeparatorFromArg(args[argi+1]);
		//		writerOptions.output_json_flatten_separator = SeparatorFromArg(args[argi+1]);
		//		argi += 2;
		//
		//	} else if args[argi] == "--io" {
		//		checkArgCount(args, argi, argc, 2);
		//		if (!lhmss_has_key(get_default_rses(), args[argi+1])) {
		//			fmt.Fprintf(os.Stderr, "%s: unrecognized I/O format \"%s\".\n",
		//				os.Args[0], args[argi+1]);
		//			os.Exit(1);
		//		}
		//		readerOptions.InputFileFormat = args[argi+1];
		//		writerOptions.OutputFileFormat = args[argi+1];
		//		argi += 2;
		//
	} else if args[argi] == "--csv" {
		readerOptions.InputFileFormat = "csv"
		writerOptions.OutputFileFormat = "csv"
		argi += 1

		//	} else if args[argi] == "--csvlite" {
		//		readerOptions.InputFileFormat = "csvlite";
		//		writerOptions.OutputFileFormat = "csvlite";
		//		argi += 1;
		//
		//	} else if args[argi] == "--tsv" {
		//		readerOptions.InputFileFormat = writerOptions.OutputFileFormat = "csv";
		//		readerOptions.IFS = "\t";
		//		writerOptions.OFS = "\t";
		//		argi += 1;
		//
		//	} else if args[argi] == "--tsvlite" || args[argi] == "-t" {
		//		readerOptions.InputFileFormat = writerOptions.OutputFileFormat = "csvlite";
		//		readerOptions.IFS = "\t";
		//		writerOptions.OFS = "\t";
		//		argi += 1;
		//
		//	} else if args[argi] == "--asv" {
		//		readerOptions.InputFileFormat = writerOptions.OutputFileFormat = "csv";
		//		readerOptions.IFS = ASV_FS;
		//		writerOptions.OFS = ASV_FS;
		//		readerOptions.IRS = ASV_RS;
		//		writerOptions.ORS = ASV_RS;
		//		argi += 1;
		//
		//	} else if args[argi] == "--asvlite" {
		//		readerOptions.InputFileFormat = writerOptions.OutputFileFormat = "csvlite";
		//		readerOptions.IFS = ASV_FS;
		//		writerOptions.OFS = ASV_FS;
		//		readerOptions.IRS = ASV_RS;
		//		writerOptions.ORS = ASV_RS;
		//		argi += 1;
		//
		//	} else if args[argi] == "--usv" {
		//		readerOptions.InputFileFormat = writerOptions.OutputFileFormat = "csv";
		//		readerOptions.IFS = USV_FS;
		//		writerOptions.OFS = USV_FS;
		//		readerOptions.IRS = USV_RS;
		//		writerOptions.ORS = USV_RS;
		//		argi += 1;
		//
		//	} else if args[argi] == "--usvlite" {
		//		readerOptions.InputFileFormat = writerOptions.OutputFileFormat = "csvlite";
		//		readerOptions.IFS = USV_FS;
		//		writerOptions.OFS = USV_FS;
		//		readerOptions.IRS = USV_RS;
		//		writerOptions.ORS = USV_RS;
		//		argi += 1;
		//
	} else if args[argi] == "--dkvp" {
		readerOptions.InputFileFormat = "dkvp"
		writerOptions.OutputFileFormat = "dkvp"
		argi += 1

	} else if args[argi] == "--json" {
		readerOptions.InputFileFormat = "json"
		writerOptions.OutputFileFormat = "json"
		argi += 1
		//	} else if args[argi] == "--jsonx" {
		//		readerOptions.InputFileFormat = "json";
		//		writerOptions.OutputFileFormat = "json";
		//		writerOptions.stack_json_output_vertically = true;
		//		argi += 1;
		//
	} else if args[argi] == "--nidx" {
		readerOptions.InputFileFormat = "nidx"
		writerOptions.OutputFileFormat = "nidx"
		argi += 1

		//	} else if args[argi] == "-T" {
		//		readerOptions.InputFileFormat = "nidx";
		//		writerOptions.OutputFileFormat = "nidx";
		//		readerOptions.IFS = "\t";
		//		writerOptions.OFS = "\t";
		//		argi += 1;
		//
		//	} else if args[argi] == "--xtab" {
		//		readerOptions.InputFileFormat = "xtab";
		//		writerOptions.OutputFileFormat = "xtab";
		//		argi += 1;
		//
		//	} else if args[argi] == "--pprint" {
		//		readerOptions.InputFileFormat        = "csvlite";
		//		readerOptions.IFS              = " ";
		//		readerOptions.allow_repeat_ifs = true;
		//		writerOptions.OutputFileFormat        = "pprint";
		//		argi += 1;
		//
		//	} else if args[argi] == "--c2t" {
		//		readerOptions.InputFileFormat = "csv";
		//		readerOptions.IRS       = "auto";
		//		writerOptions.OutputFileFormat = "csv";
		//		writerOptions.ORS       = "auto";
		//		writerOptions.OFS       = "\t";
		//		argi += 1;
		//	} else if args[argi] == "--c2d" {
		//		readerOptions.InputFileFormat = "csv";
		//		readerOptions.IRS       = "auto";
		//		writerOptions.OutputFileFormat = "dkvp";
		//		argi += 1;
		//	} else if args[argi] == "--c2n" {
		//		readerOptions.InputFileFormat = "csv";
		//		readerOptions.IRS       = "auto";
		//		writerOptions.OutputFileFormat = "nidx";
		//		argi += 1;
		//	} else if args[argi] == "--c2j" {
		//		readerOptions.InputFileFormat = "csv";
		//		readerOptions.IRS       = "auto";
		//		writerOptions.OutputFileFormat = "json";
		//		argi += 1;
		//	} else if args[argi] == "--c2p" {
		//		readerOptions.InputFileFormat = "csv";
		//		readerOptions.IRS       = "auto";
		//		writerOptions.OutputFileFormat = "pprint";
		//		argi += 1;
		//	} else if args[argi] == "--c2x" {
		//		readerOptions.InputFileFormat = "csv";
		//		readerOptions.IRS       = "auto";
		//		writerOptions.OutputFileFormat = "xtab";
		//		argi += 1;
		//	} else if args[argi] == "--c2m" {
		//		readerOptions.InputFileFormat = "csv";
		//		readerOptions.IRS       = "auto";
		//		writerOptions.OutputFileFormat = "markdown";
		//		argi += 1;
		//
		//	} else if args[argi] == "--t2c" {
		//		readerOptions.InputFileFormat = "csv";
		//		readerOptions.IFS       = "\t";
		//		readerOptions.IRS       = "auto";
		//		writerOptions.OutputFileFormat = "csv";
		//		writerOptions.ORS       = "auto";
		//		argi += 1;
		//	} else if args[argi] == "--t2d" {
		//		readerOptions.InputFileFormat = "csv";
		//		readerOptions.IFS       = "\t";
		//		readerOptions.IRS       = "auto";
		//		writerOptions.OutputFileFormat = "dkvp";
		//		argi += 1;
		//	} else if args[argi] == "--t2n" {
		//		readerOptions.InputFileFormat = "csv";
		//		readerOptions.IFS       = "\t";
		//		readerOptions.IRS       = "auto";
		//		writerOptions.OutputFileFormat = "nidx";
		//		argi += 1;
		//	} else if args[argi] == "--t2j" {
		//		readerOptions.InputFileFormat = "csv";
		//		readerOptions.IFS       = "\t";
		//		readerOptions.IRS       = "auto";
		//		writerOptions.OutputFileFormat = "json";
		//		argi += 1;
		//	} else if args[argi] == "--t2p" {
		//		readerOptions.InputFileFormat = "csv";
		//		readerOptions.IFS       = "\t";
		//		readerOptions.IRS       = "auto";
		//		writerOptions.OutputFileFormat = "pprint";
		//		argi += 1;
		//	} else if args[argi] == "--t2x" {
		//		readerOptions.InputFileFormat = "csv";
		//		readerOptions.IFS       = "\t";
		//		readerOptions.IRS       = "auto";
		//		writerOptions.OutputFileFormat = "xtab";
		//		argi += 1;
		//	} else if args[argi] == "--t2m" {
		//		readerOptions.InputFileFormat = "csv";
		//		readerOptions.IFS       = "\t";
		//		readerOptions.IRS       = "auto";
		//		writerOptions.OutputFileFormat = "markdown";
		//		argi += 1;
		//
		//	} else if args[argi] == "--d2c" {
		//		readerOptions.InputFileFormat = "dkvp";
		//		writerOptions.OutputFileFormat = "csv";
		//		writerOptions.ORS       = "auto";
		//		argi += 1;
		//	} else if args[argi] == "--d2t" {
		//		readerOptions.InputFileFormat = "dkvp";
		//		writerOptions.OutputFileFormat = "csv";
		//		writerOptions.ORS       = "auto";
		//		writerOptions.OFS       = "\t";
		//		argi += 1;
		//	} else if args[argi] == "--d2n" {
		//		readerOptions.InputFileFormat = "dkvp";
		//		writerOptions.OutputFileFormat = "nidx";
		//		argi += 1;
		//	} else if args[argi] == "--d2j" {
		//		readerOptions.InputFileFormat = "dkvp";
		//		writerOptions.OutputFileFormat = "json";
		//		argi += 1;
		//	} else if args[argi] == "--d2p" {
		//		readerOptions.InputFileFormat = "dkvp";
		//		writerOptions.OutputFileFormat = "pprint";
		//		argi += 1;
		//	} else if args[argi] == "--d2x" {
		//		readerOptions.InputFileFormat = "dkvp";
		//		writerOptions.OutputFileFormat = "xtab";
		//		argi += 1;
		//	} else if args[argi] == "--d2m" {
		//		readerOptions.InputFileFormat = "dkvp";
		//		writerOptions.OutputFileFormat = "markdown";
		//		argi += 1;
		//
		//	} else if args[argi] == "--n2c" {
		//		readerOptions.InputFileFormat = "nidx";
		//		writerOptions.OutputFileFormat = "csv";
		//		writerOptions.ORS       = "auto";
		//		argi += 1;
		//	} else if args[argi] == "--n2t" {
		//		readerOptions.InputFileFormat = "nidx";
		//		writerOptions.OutputFileFormat = "csv";
		//		writerOptions.ORS       = "auto";
		//		writerOptions.OFS       = "\t";
		//		argi += 1;
		//	} else if args[argi] == "--n2d" {
		//		readerOptions.InputFileFormat = "nidx";
		//		writerOptions.OutputFileFormat = "dkvp";
		//		argi += 1;
		//	} else if args[argi] == "--n2j" {
		//		readerOptions.InputFileFormat = "nidx";
		//		writerOptions.OutputFileFormat = "json";
		//		argi += 1;
		//	} else if args[argi] == "--n2p" {
		//		readerOptions.InputFileFormat = "nidx";
		//		writerOptions.OutputFileFormat = "pprint";
		//		argi += 1;
		//	} else if args[argi] == "--n2x" {
		//		readerOptions.InputFileFormat = "nidx";
		//		writerOptions.OutputFileFormat = "xtab";
		//		argi += 1;
		//	} else if args[argi] == "--n2m" {
		//		readerOptions.InputFileFormat = "nidx";
		//		writerOptions.OutputFileFormat = "markdown";
		//		argi += 1;
		//
		//	} else if args[argi] == "--j2c" {
		//		readerOptions.InputFileFormat = "json";
		//		writerOptions.OutputFileFormat = "csv";
		//		writerOptions.ORS       = "auto";
		//		argi += 1;
		//	} else if args[argi] == "--j2t" {
		//		readerOptions.InputFileFormat = "json";
		//		writerOptions.OutputFileFormat = "csv";
		//		writerOptions.ORS       = "auto";
		//		writerOptions.OFS       = "\t";
		//		argi += 1;
		//	} else if args[argi] == "--j2d" {
		//		readerOptions.InputFileFormat = "json";
		//		writerOptions.OutputFileFormat = "dkvp";
		//		argi += 1;
		//	} else if args[argi] == "--j2n" {
		//		readerOptions.InputFileFormat = "json";
		//		writerOptions.OutputFileFormat = "nidx";
		//		argi += 1;
		//	} else if args[argi] == "--j2p" {
		//		readerOptions.InputFileFormat = "json";
		//		writerOptions.OutputFileFormat = "pprint";
		//		argi += 1;
		//	} else if args[argi] == "--j2x" {
		//		readerOptions.InputFileFormat = "json";
		//		writerOptions.OutputFileFormat = "xtab";
		//		argi += 1;
		//	} else if args[argi] == "--j2m" {
		//		readerOptions.InputFileFormat = "json";
		//		writerOptions.OutputFileFormat = "markdown";
		//		argi += 1;
		//
		//	} else if args[argi] == "--p2c" {
		//		readerOptions.InputFileFormat        = "csvlite";
		//		readerOptions.IFS              = " ";
		//		readerOptions.allow_repeat_ifs = true;
		//		writerOptions.OutputFileFormat        = "csv";
		//		writerOptions.ORS              = "auto";
		//		argi += 1;
		//	} else if args[argi] == "--p2t" {
		//		readerOptions.InputFileFormat        = "csvlite";
		//		readerOptions.IFS              = " ";
		//		readerOptions.allow_repeat_ifs = true;
		//		writerOptions.OutputFileFormat        = "csv";
		//		writerOptions.ORS              = "auto";
		//		writerOptions.OFS              = "\t";
		//		argi += 1;
		//	} else if args[argi] == "--p2d" {
		//		readerOptions.InputFileFormat        = "csvlite";
		//		readerOptions.IFS              = " ";
		//		readerOptions.allow_repeat_ifs = true;
		//		writerOptions.OutputFileFormat        = "dkvp";
		//		argi += 1;
		//	} else if args[argi] == "--p2n" {
		//		readerOptions.InputFileFormat        = "csvlite";
		//		readerOptions.IFS              = " ";
		//		readerOptions.allow_repeat_ifs = true;
		//		writerOptions.OutputFileFormat        = "nidx";
		//		argi += 1;
		//	} else if args[argi] == "--p2j" {
		//		readerOptions.InputFileFormat        = "csvlite";
		//		readerOptions.IFS              = " ";
		//		readerOptions.allow_repeat_ifs = true;
		//		writerOptions.OutputFileFormat = "json";
		//		argi += 1;
		//	} else if args[argi] == "--p2x" {
		//		readerOptions.InputFileFormat        = "csvlite";
		//		readerOptions.IFS              = " ";
		//		readerOptions.allow_repeat_ifs = true;
		//		writerOptions.OutputFileFormat        = "xtab";
		//		argi += 1;
		//	} else if args[argi] == "--p2m" {
		//		readerOptions.InputFileFormat        = "csvlite";
		//		readerOptions.IFS              = " ";
		//		readerOptions.allow_repeat_ifs = true;
		//		writerOptions.OutputFileFormat        = "markdown";
		//		argi += 1;
		//
		//	} else if args[argi] == "--x2c" {
		//		readerOptions.InputFileFormat = "xtab";
		//		writerOptions.OutputFileFormat = "csv";
		//		writerOptions.ORS       = "auto";
		//		argi += 1;
		//	} else if args[argi] == "--x2t" {
		//		readerOptions.InputFileFormat = "xtab";
		//		writerOptions.OutputFileFormat = "csv";
		//		writerOptions.ORS       = "auto";
		//		writerOptions.OFS       = "\t";
		//		argi += 1;
		//	} else if args[argi] == "--x2d" {
		//		readerOptions.InputFileFormat = "xtab";
		//		writerOptions.OutputFileFormat = "dkvp";
		//		argi += 1;
		//	} else if args[argi] == "--x2n" {
		//		readerOptions.InputFileFormat = "xtab";
		//		writerOptions.OutputFileFormat = "nidx";
		//		argi += 1;
		//	} else if args[argi] == "--x2j" {
		//		readerOptions.InputFileFormat = "xtab";
		//		writerOptions.OutputFileFormat = "json";
		//		argi += 1;
		//	} else if args[argi] == "--x2p" {
		//		readerOptions.InputFileFormat = "xtab";
		//		writerOptions.OutputFileFormat = "pprint";
		//		argi += 1;
		//	} else if args[argi] == "--x2m" {
		//		readerOptions.InputFileFormat = "xtab";
		//		writerOptions.OutputFileFormat = "markdown";
		//		argi += 1;
		//
		//	} else if args[argi] == "-N" {
		//		readerOptions.use_implicit_csv_header = true;
		//		writerOptions.headerless_csv_output = true;
		//		argi += 1;
	}
	*pargi = argi
	return argi != oargi
}

// Returns true if the current flag was handled.
func handleMiscOptions(args []string, argc int, pargi *int, options *clitypes.TOptions) bool {
	argi := *pargi
	oargi := argi
	//
	//	if args[argi] == "-I" {
	//		options.do_in_place = true;
	//		argi += 1;
	//
	//	} else if args[argi] == "-n" {
	//		options.no_input = true;
	//		argi += 1;
	//
	//	} else if args[argi] == "--from" {
	//		checkArgCount(args, argi, argc, 2);
	//		slls_append(options.filenames, args[argi+1], NO_FREE);
	//		argi += 2;
	//
	//	} else if args[argi] == "--ofmt" {
	//		checkArgCount(args, argi, argc, 2);
	//		options.ofmt = args[argi+1];
	//		argi += 2;
	//
	//	} else if args[argi] == "--nr-progress-mod" {
	//		checkArgCount(args, argi, argc, 2);
	//		if (sscanf(args[argi+1], "%lld", &options.nr_progress_mod) != 1) {
	//			fmt.Fprintf(os.Stderr,
	//				"%s: --nr-progress-mod argument must be a positive integer; got \"%s\".\n",
	//				os.Args[0], args[argi+1]);
	//			mainUsageShort()
	//			os.Exit(1);
	//		}
	//		if (options.nr_progress_mod <= 0) {
	//			fmt.Fprintf(os.Stderr,
	//				"%s: --nr-progress-mod argument must be a positive integer; got \"%s\".\n",
	//				os.Args[0], args[argi+1]);
	//			mainUsageShort()
	//			os.Exit(1);
	//		}
	//		argi += 2;
	//
	//	} else if args[argi] == "--seed" {
	//		checkArgCount(args, argi, argc, 2);
	//		if (sscanf(args[argi+1], "0x%x", &options.rand_seed) == 1) {
	//			options.have_rand_seed = true;
	//		} else if (sscanf(args[argi+1], "%u", &options.rand_seed) == 1) {
	//			options.have_rand_seed = true;
	//		} else {
	//			fmt.Fprintf(os.Stderr,
	//				"%s: --seed argument must be a decimal or hexadecimal integer; got \"%s\".\n",
	//				os.Args[0], args[argi+1]);
	//			mainUsageShort()
	//			os.Exit(1);
	//		}
	//		argi += 2;
	//
	//	}
	*pargi = argi
	return argi != oargi
}

// ----------------------------------------------------------------
//static char* lhmss_get_or_die(lhmss_t* pmap, char* key) {
//	char* value = lhmss_get(pmap, key);
//	MLR_INTERNAL_CODING_ERROR_IF(value == nil);
//	return value;
//}

// ----------------------------------------------------------------
//static int lhmsll_get_or_die(lhmsll_t* pmap, char* key) {
//	MLR_INTERNAL_CODING_ERROR_UNLESS(lhmsll_has_key(pmap, key));
//	return lhmsll_get(pmap, key);
//}

//	cpuprofile := flag.String("cpuprofile", "", "Write CPU profile to `file`")
//// ----------------------------------------------------------------
//func maybeProfile(cpuprofile *string) {
//	// to do: move to method
//	// go tool pprof mlr foo.prof
//	//   top10
//	if *cpuprofile != "" {
//		f, err := os.Create(*cpuprofile)
//		if err != nil {
//			fmt.Fprintln(os.Stderr, os.Args[0], ": ", "Could not start CPU profile: ", err)
//		}
//		defer f.Close()
//		if err := pprof.StartCPUProfile(f); err != nil {
//			fmt.Fprintln(os.Stderr, os.Args[0], ": ", "Could not start CPU profile: ", err)
//		}
//		defer pprof.StopCPUProfile()
//	}