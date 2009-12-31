<?php
/**
 *  Class: Controller
 *
 *  The Controller class should be the parent class of all of your Controller sub classes
 *  that contain the business logic of your application (render a blog post, log a user in,
 *  delete something and redirect, etc).
 *
 */

class Controller
{
	protected $layout = false;
	protected $layout_vars = array();

	/*
	   Method: execute
	   Method use by the controller to validate and execute the action
	*/
	public function execute($action, $params)
	{
		// it's a private method of the class or action is not a method of the class
		if (substr($action, 0, 1) == '_' or !method_exists($this, $action))
			throw new Exception("Action '{$action}' is not valid!");

		call_user_func_array(array($this, $action), $params);
	}

	/*
	   Method: setLayout
	   Set the layout file

	   the file will use the same extention as the view
	*/
	public function setLayout($layout)
	{
		$this->layout = $layout;
	}

	/*
	   Method: assign
	   Assign specific variable to the layout
	*/
	public function assignToLayout($var, $value)
	{
		if (is_array($var))
			array_merge($this->layout_vars, $var);
		else
			$this->layout_vars[$var] = $value;
	}

	/*
	   Method: render
	   Render the view and the layout if setted and return it as a string
	*/
	public function render($view, $vars=array())
	{
		if ($this->layout)
		{
			$this->layout_vars['content_for_layout'] = new View($view, $vars);
			return new View('../layouts/'.$this->layout, $this->layout_vars);
		}
		else return new View($view, $vars);
	}

	/*
	   Method: display
	   Display the view and the layout if setted
	*/
	public function display($view, $vars=array(), $exit=true)
	{
		Observer::notify('system.display');

		echo $this->render($view, $vars);

		if ($exit)
		{
			Observer::notify('system.shutdown');
			exit;
		}
	}

} // end Controller class