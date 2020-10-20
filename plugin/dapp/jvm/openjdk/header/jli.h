#ifndef _JLI_H_
#define _JLI_H_



int JLI_Exec_Contract(int argc, char **argv, char **exceptionInfo, int jobType, char *jvmGo);
int JLI_Detroy_JVM();
int JLI_Create_JVM(const char *jdkPath);


/* utility functions */
extern int GetPtrSize();
extern void SetPtr(char **ptr, char *value, int index);
extern void FreeArgv(int argc, char **argv);
extern char ** GetNil2dPtr();
extern void * GetVoidPtr(char *voidPtr);

#endif /* _JAVA_H_ */
