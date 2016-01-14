from threading import Thread
#hshshshsh
i = 0

class MyThread(Thread):
	def __init__(self, val):
		''' Constructor. '''
		Thread.__init__(self)
		self.val = val
	def run(self):
		global i
		if self.val == 0:
			for j in range(0, 1000000):
				i +=1
		if self.val == 1:
			for j in range(0, 1000000):
				i -=1

def main():
	myThreadOb1 = MyThread(0)
	myThreadOb1.setName('Thread 1')
	myThreadOb2 = MyThread(1)
	myThreadOb2.setName('Thread 2')
	# Start running the threads!
	myThreadOb1.start()
	myThreadOb2.start()
	myThreadOb1.join()
	myThreadOb2.join()
	print i

main()
