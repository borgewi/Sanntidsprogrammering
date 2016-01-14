#include <pthread.h>
#include <stdlib.h>
#include <stdio.h>
#include <string.h>


int i = 0;

void* func1(){
	int j;
	for (j = 0; j<1000000;j++){
		i++;
	} 
}

void* func2(){
	int j;
	for (j = 0; j<1000000;j++){
		i--;
	} 
}

int main(void){
	pthread_t hilo1, hilo2;
	pthread_create(&hilo1,NULL, &func1, NULL);
	pthread_create(&hilo2,NULL, &func2, NULL);

	printf("the main thread continues with its execution\n");

	pthread_join(hilo1,NULL);
	pthread_join(hilo2, NULL);

	printf("the main thread finished");
	printf("%d",i);
	printf("\n");

	return(0);

}
