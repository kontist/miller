#ifndef CONTEXT_H
#define CONTEXT_H

// File-level context for Miller's NR, FNR, FILENAME, and FILENUM variables, as
// well as for error messages
typedef struct _context_t {
	long long nr;
	long long fnr;
	int       filenum;
	char*     filename;
	int       force_eof; // e.g. mlr head

	char*     ips;
	char*     ifs;
	char*     irs;
	char*     ops;
	char*     ofs;
	char*     ors;

	// For autodetect between LF and CRLF line endings: the lrec reader
	// (when in --auto mode) can place here the line ending which was
	// encountered on the first record read.
	char*     auto_line_term;
	int       auto_line_term_detected;
} context_t;

void context_init_from_first_file_name(context_t* pctx, char* first_file_name);
void context_init_from_opts(context_t* pctx, void* pvopts);

void context_set_autodetected_crlf(context_t* pctx);
void context_set_autodetected_lf(context_t* pctx);
void context_set_autodetected_line_term(context_t* pctx, char* line_term);

void context_print(context_t* pctx, char* indent);

#endif // CONTEXT_H
