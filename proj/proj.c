#include "proj_go.h"

void trans(PJ *pj, PJ_DIRECTION direction, double u1, double v1, double w1, double t1, double *u2, double *v2, double *w2, double *t2) {
    PJ_COORD
	co1,
	co2;
    co1.uvwt.u = u1;
    co1.uvwt.v = v1;
    co1.uvwt.w = w1;
    co1.uvwt.t = t1;
    co2 = proj_trans(pj, direction, co1);
    *u2 = co2.uvwt.u;
    *v2 = co2.uvwt.v;
    *w2 = co2.uvwt.w;
    *t2 = co2.uvwt.t;    
}

PJ_COORD uvwt(double u, double v, double w, double t) {
    PJ_COORD
	c;
    c.uvwt.u = u;
    c.uvwt.v = v;
    c.uvwt.w = w;
    c.uvwt.t = t;
    return c;
}
