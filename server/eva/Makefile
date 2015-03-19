
CC=cc

allfive: allfive.c pokerlib.o
	${CC} ${CFLAGS} allfive.c pokerlib.o -o allfive

pokerlib.o: pokerlib.c arrays.h
	${CC} -c ${CFLAGS} pokerlib.c -o pokerlib.o