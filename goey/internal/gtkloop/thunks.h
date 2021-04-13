#ifndef GOEY_GTKLOOP_THUNKS_H
#define GOEY_GTKLOOP_THUNKS_H

extern void loopInit(void);
extern void loopRun(void);
extern void loopStop(void);
extern void loopMainContextInvoke(void);
extern void loopIdleAdd(void);

#endif