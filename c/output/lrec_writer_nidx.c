#include <stdlib.h>
#include "lib/mlrutil.h"
#include "output/lrec_writers.h"

typedef struct _lrec_writer_nidx_state_t {
	char* ors;
	char* ofs;
} lrec_writer_nidx_state_t;

static void lrec_writer_nidx_free(lrec_writer_t* pwriter);
static void lrec_writer_nidx_process(void* pvstate, FILE* output_stream, lrec_t* prec);

// ----------------------------------------------------------------
lrec_writer_t* lrec_writer_nidx_alloc(char* ors, char* ofs) {
	lrec_writer_t* plrec_writer = mlr_malloc_or_die(sizeof(lrec_writer_t));

	lrec_writer_nidx_state_t* pstate = mlr_malloc_or_die(sizeof(lrec_writer_nidx_state_t));
	pstate->ors = ors;
	pstate->ofs = ofs;

	plrec_writer->pvstate       = (void*)pstate;
	plrec_writer->pprocess_func = lrec_writer_nidx_process;
	plrec_writer->pfree_func    = lrec_writer_nidx_free;

	return plrec_writer;
}

static void lrec_writer_nidx_free(lrec_writer_t* pwriter) {
	free(pwriter->pvstate);
	free(pwriter);
}

// ----------------------------------------------------------------
static void lrec_writer_nidx_process(void* pvstate, FILE* output_stream, lrec_t* prec) {
	if (prec == NULL)
		return;
	lrec_writer_nidx_state_t* pstate = pvstate;
	char* ors = pstate->ors;
	char* ofs = pstate->ofs;

	int nf = 0;
	for (lrece_t* pe = prec->phead; pe != NULL; pe = pe->pnext) {
		if (nf > 0)
			fputs(ofs, output_stream);
		fputs(pe->value, output_stream);
		nf++;
	}
	fputs(ors, output_stream);
	lrec_free(prec); // end of baton-pass
}
