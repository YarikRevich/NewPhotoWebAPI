package common
/**
	IT IS A GENERAL TUTORIAL OF A GOOD LOOKING MODEL...

	Naming:
	  
		"{*TYPE_OF_REQUEST}{**TYPE_OF_MODEL}{NAME_OF_HANDLER}Model"

		* - "POST", "GET", "PUT" ...
		** - "Request" of "Response"

	Struct:

		Request -> 
		
		 {
			 "data": {...}
		 }

		Response ->

		 {
			 "result": {...} (Optional)
			 "service": {
				 "ok": ...
				 "message": ... (Optional)
			 }
		 }
**/