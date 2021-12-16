# EstimadoresDFSA
FLAGS:

-eom-lee

    enable Eom-Lee estimator
    
    
-frame-two
    
    limit frame size to a power of two number
    
    
-inc-tags int

    number of tags to increment by each step
    
    
-iv-ii
   
   enable IV-II estimator
    
    
-lower-bound
   
   enable lower bound estimator
    
    
-max-tags int
  
  number of maximum tags to simulate
    
    
-replay-step int
  
  number of iterations by each step
    
    
-shoute
  
  enable shoute estimator
  
  
-start-frame int
  
  initial value of frame size
  
  
-start-tags int
    initial number of tags to start simulation
    
How to run:
  Go to the projects folder and do:
  `go run main.go  [TAGS]`
  To run you must select at least one estimator, and set -start-tags and -start-frame
