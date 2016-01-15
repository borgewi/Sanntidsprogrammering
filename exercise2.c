#include <pthread.h>
#include <stdlib.h>
#include <stdio.h>
#include <string.h>


pthread_mutex_t lock;

int i = 0;

void* func1(){
	pthread_mutex_lock(&lock);
	int j;
	for (j = 0; j<1000000;j++){
		i++;
	}
	 pthread_mutex_unlock(&lock);
}

void* func2(){
	pthread_mutex_lock(&lock);
	int j;
	for (j = 0; j<1000000;j++){
		i--;
	}
	pthread_mutex_unlock(&lock);
}

int main(void){
	pthread_t hilo1, hilo2;
	pthread_create(&hilo1,NULL, &func1, NULL);
	pthread_create(&hilo2,NULL, &func2, NULL);

	printf("the main thread continues with its execution\n");

	pthread_join(hilo1,NULL);
	pthread_join(hilo2, NULL);
	pthread_mutex_destroy(&lock);

	printf("the main thread finished");
	printf("%d",i);
	printf("\n");

	return(0);

}
