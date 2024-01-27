<?php
/*
   Class: View

   The template object takes a valid path to a template file as the only
   argument required, but you can add directly all properties to the template,
   which become available as local vars in the template file. You can then call
   display() to get the output of the template, or just call print on the
   template directly thanks to PHP 5's __toString magic method.

   > echo new View('my_template', array(
   >     'title' => 'My Title',
   >     'body' => 'My body content'
   > ));

   my_template.php might look like this:
   > <html>
   > <head><title><?php echo $title;?></title></head>
   > <body>
   >   <h1><?php echo $title;?></h1>
   >   <p><?php echo $body;?></p>
   > </body>
   > </html>
*/

class View
{
	private $file;           // path of template file
	private $vars = array(); // array of template variables

	public function __construct($file, $vars=false)
	{
		$this->file = APPPATH.'/views/'.$file.VIEW_SUFFIX;

		if ( ! file_exists($this->file))
			throw new Exception("View '{$file}' not found!");

		if ($vars !== false)
			$this->vars = $vars;
	}

	/*
	   Method: assign
	   Assign specific variable to the template
	*/
	public function assign($name, $value=null)
	{
		if (is_array($name))
			array_merge($this->vars, $name);
		else
			$this->vars[$name] = $value;
	}

	/*
	   Method: render
	   Render the template and return it as string
	*/
	public function render()
	{
		ob_start();

		extract($this->vars, EXTR_SKIP);
		include $this->file;

		$content = ob_get_clean();
		return $content;
	}

	/*
	   Method: display
	   Display the template
	*/
	public function display()
	{
		echo $this->render();
	}

	public function __toString() { return $this->render(); }

} // end View class