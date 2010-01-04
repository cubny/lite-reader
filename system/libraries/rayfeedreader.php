<?php

    /**
     * RayFeedReader
     *
     * SimpleXML based feed reader class. A very specific feed reader designed
     * to working with no or little configurations.
     *
     * - This class can read an rss feed into array from a given url.
     * - Also can render html widget with plugable RayFeedWidget Class.
     * 
     * Supports following feed types:
     *  - RSS 0.92
     *  - RSS 2.0
     *  - RDF
     *  - Atom
     * 
     * Public Methods
     * - getInstance()
     * - setOptions()
     * - parse()
     * - getData()
     * - getType()
     * - widget()
     * 
     * 
     * Configuration Options
     *  - array
     *      - url: (string)
     *          - feed url
     *
     *      - httpClient: (string)
     *          - default php (native file_get_contents)
     *          - value php or rayHttp or SimpleXML
     *
     *      - type: ([optional] string
     *          - auto detect
     *          - value rss or rss2 or rdf or atom
     *
     *      - widget: ([optional] string)
     *          - feed widget class name for rendering html
     * 
     *      - rayHttp: (array)
     *          - only if httpClient is set to rayHttp
     *          - rayHttp Options if you want to modify rayHttml CURL options
     *          - generally not required.
     *
     *      - followCdata: (boolean)
     *          - default true
     *          - partially implemented
     *
     *
     *
     * @version 1.01
     * @author Md. Rayhan Chowdhury
     * @package rayFeedReader
     * @license GPL
     */
    Class RayFeedReader{

        /**
         * Self Instance for Singleton Pattern
         *
         * @var object
         * @access protected
         */
        static private  $__instance;

        /**
         * Instance of Parser Class.
         *
         * @var object Parser Class
         * @access protected
         */
        Protected       $_Parser;

        /**
         * Feed Url
         *
         * @var string feed url
         * @access protected
         */
        protected       $_url;
        
        /**
         * XML Feed content
         *
         * @var string [optional]
         * @access protected
         */
        protected		$_xml;

        /**
         * List if error messages
         *
         * @var array
         * @access private
         */
        private         $__errors = array();

        /**
         * Runtime Options for reader
         *
         * @var array
         * @access protected
         */
        protected       $_options = array('rayHttp' => array());

        /**
         * Type of feed to be parsed.
         *
         * @var string
         * @access protected
         */
        protected       $_type = null;

        /**
         * Detect feed type automatically.
         *
         * @var boolean
         * @access protected
         */
        protected         $_autoDetectFeedType = true;

        /**
         * HttpClient to be used for loading feed content.
         *
         *  - default SimpleXML
         *
         * @var string 'SimpleXML' or 'rayHttp'
         * @access protected
         */
        protected       $_httpClient = "php";
        

        /**
         * Follow CDATA elements.
         *
         *  - default true
         *
         * @var boolean
         * @access protected
         */
        protected       $_followCdata = true;

        /**
         * Widget Class Name
         *
         * @var string
         * @access protected
         */
        protected       $_widget;

        /**
         * Parsed result data
         * 
         * @var array
         * @access protected
         */
        protected       $_content;

        /**
         * Class construct
         *
         * @param array $options
         */
        function __construct($options = array()) {
            $this->setOptions($options);
        }

        /**
         * Get Error occured during parsing a fead.
         *
         * @param boolean $clean, clean all errors
         * @return array
         * @access public
         */
        function getErrors($clean = true) {
            if ($clean) {
                $errors = $this->__errors;
                unset($this->__errors);
                return $errors;
            }else {
                return $this->__errors;
            }
        }

        /**
         * Get Instance of the class.
         *
         * @param array $options
         * @return object self instance.
         * @access public
         * @static
         */
        static function &getInstance($options = array()) {
            if (is_null(self::$__instance)) {
                self::$__instance = new self($options);
            }
            return self::$__instance;
        }

        /**
         * Set Options for the class
         * 
         * 
         * @param array $options
         * @return object self instance
         * @access public
         */
        function &setOptions($options) {
            if (!empty($options['url'])) {
                $this->_url = $options['url'];
            } elseif (is_null($options['url'])) {
            	$this->_url = null;
            }
            
            if (!empty($options['xml'])) {
            	$this->_xml = $options['xml'];
            	
            }

            if (isset($options['type'])) {
                /**
                 * Set/Unset feed auto detection...
                 */
                if (!empty($options['type'])) {
                    $this->_autoDetectFeedType = false;
                    $this->_type = $options['type'];
                } else {
                    $this->_autoDetectFeedType = true;
                }
            }


            if (!empty($options['_type'])) {
                    $this->_type = $options['_type'];
            }

            if (!empty($options['httpClient'])) {
                $this->_httpClient = $options['httpClient'];
            }

            if (isset($options['followCdata'])) {
                $this->_followCdata = $options['followCdata'];
            }

            if (!empty($options['widget'])) {
                $this->_widget = $options['widget'];
            }

            $this->_options = array_merge($this->_options, $options);

            return $this;
        }

        /**
         * Parse feed contents into an array and return self object
         * 
         * @return object self instance
         * @access public
         */
        function &parse() {
            /**
             * reset contents.
             */
            $this->_content = null;
			
		        try{
		            libxml_use_internal_errors(true);
		            libxml_clear_errors(true);
                  
                  $cdata = LIBXML_NOCDATA;
                  if (!$this->_followCdata) {
                     $cdata = null;                  
                  }
                  
		            if (!is_null($this->_url)) {
				        /**
				         * Get/load content
				         */
				         switch ($this->_httpClient) {
				             case 'php':
				                 $content = @file_get_contents($this->_url);
				                 if (!empty($content)) {
				                    $content = new SimpleXMLElement($content); //, $cdata
				                 }
				                 break;

				             case 'SimpleXML':
				                 $content = new SimpleXMLElement($this->_url, $cdata, true);
				                 break;

				             case 'rayHttp':
				                 $content = RayHttp::getInstance()->setOptions($this->_options['rayHttp'])->get($this->_url);
				                 if (!empty($content)) {
				                    $content = new SimpleXMLElement($content, $cdata);
				                 }
				                 break;
				         }
		             } else {
		             	// If xml feed stream is passed
		             	if (!empty($this->_xml)) {
		                    $content = new SimpleXMLElement($this->_xml, $cdata);
		             	}
		             }

		        } catch (Exception $e) {
		            $this->__errors[] = $e->getMessage();
		            return $this;
		        }
            
            if (empty($content)) {
                $this->__errors[] = 'Feed XML is either invalid or broken.';
                return $this;
            }

            /**
             * Detect Feed Type
             */
             if ($this->_autoDetectFeedType) {
                    
                    switch ($content->getName()) {
                        case 'rss':
                            foreach ($content->attributes() as $attribute) {
                                if ($attribute->getName() == 'version') {
                                    if ('2.0' == $attribute) {
                                        self::setOptions(array('_type' => 'rss2'));
                                    } elseif (in_array($attribute, array('0.92', '0.91'))) {
                                        self::setOptions(array('_type' => 'rss'));
                                    }
                                }
                            }
                            break;
                            
                        case 'RDF':
                              self::setOptions(array('_type' => 'rdf'));
                            break;

                        case 'feed':                            
                            self::setOptions(array('_type' => 'atom'));
                            
                            break;
                    }
                 
             }
             
             if (!in_array($this->_type, array('rss', 'rss2', 'rdf', 'atom'))) {                 
                  $this->__errors[] = "Feed type is either invalid or not supported.";                  
                  return $this;
             }

             
            /**
             * Parse Feed Content
             */
            switch ($this->_type) {
                case 'rss':
                    $content = $this->parseRss($content);
                    break;

                case 'rdf':
                    $content = $this->parseRdf($content);
                    break;

                case 'rss2':
                    $content = $this->parseRss2($content);
                    break;

                case 'atom':
                    $content = $this->parseAtom($content);
                    break;
            }

            if (empty($content)) {
                $this->__errors[] = "No content found.";
                return $this;
            }

            $this->_content = $content;
            return $this;

        }

        /**
         * Get Array of Parsed XML feed data.
         *
         * @return array parsed feed content.
         * @access public
         */
        function getData() {
            return $this->_content;
        }

        /**
         * Get last parsed feed type
         *
         * @return string return false on failure
         * @access public
         */
        function getType() {
        	if (!empty($this->_content['type'])) {
            	return $this->_content['type'];
            }
            return false;
        }

        /**
         * Return html widget based rendered by widget class
         *
         *
         * @param array $options for html widget class
         * @return string html widget
         * @access public
         */
        function widget($options = array('widget' => 'brief')) {
            if (!empty($this->_widget) && !empty($this->_content)) {
                $Widget = new $this->_widget;
                
                return $Widget->widget($this->_content, $options);
                
             } else {
                 return false;
             }
        }
        
        /**
         * Parse feed xml into an array.
         *
         * @param object $feedXml SimpleXMLElementObject
         * @return array feed content
         * @access public
         */
        function parseRss($feedXml) {
            $data = array();

            $data['title'] = $feedXml->channel->title . '';
            $data['link'] = $feedXml->channel->link . '';
            $data['description'] = $feedXml->channel->description . '';
            $data['parser'] = __CLASS__;
            $data['type'] = 'rss';
            if (!empty($feedXml->channel->item)) {
                foreach ($feedXml->channel->item as $item) {
                    $data['items'][] = array(
                                            'title' =>  $item->title . '',
                                            'link' =>   $item->link . '',
                                            'description' => $item->description . '',
                                        );
                }
            }
            return $data;
        }

         /**
         * Parse RDF feed xml into an array.
         *
         * @param object $feedXml SimpleXMLElementObject
         * @return array feed content
         * @access public
         */
        function parseRdf($feedXml) {
            $data = array();

            $data['title'] = $feedXml->channel->title . '';
            $data['link'] = $feedXml->channel->link . '';
            $data['description'] = $feedXml->channel->description . '';
            $data['parser'] = __CLASS__;
            $data['type'] = 'rdf';

            if (!empty($feedXml->item)) {
                foreach ($feedXml->item as $item) {
                    $data['items'][] = array(
                                            'title' =>  $item->title . '',
                                            'link' =>   $item->link . '',
                                            'description' => $item->description . '',
                                        );
                }
            }
            
            return $data;
        }

        
        /**
         * Parse feed xml into an array.
         *
         * @param object $feedXml SimpleXMLElementObject
         * @return array feed content
         * @access public
         */
        function parseRss2($feedXml) {
            $data = array();

            $data['title'] = $feedXml->channel->title . '';
            $data['link'] = $feedXml->channel->link . '';
            $data['description'] = $feedXml->channel->description . '';
            $data['parser'] = __CLASS__;
            $data['type'] = 'rss2';

            $namespaces = $feedXml->getNamespaces(true);
            foreach ($namespaces as $namespace => $namespaceValue) {
                $feedXml->registerXPathNamespace($namespace, $namespaceValue);
            }
            if (!empty($feedXml->channel->item)) {
                foreach ($feedXml->channel->item as $item) {
                    $categories = array();
                    foreach ($item->children() as $child) {
                        if ($child->getName() == 'category') {
                            $categories[] = (string) $child;
                        }
                    }

                    $author = null;
                    if (!empty($namespaces['dc']) && $creator = $item->xpath('dc:creator')) {
                        $author = (string) $creator[0];
                    }

                    $content = null;
                    if (!empty($namespaces['encoded']) && $encoded = $item->xpath('content:encoded')) {
                        $content = (string) $encoded[0];
                    }

                    $data['items'][] = array(
                                            'title' =>  $item->title . '',
                                            'link' =>   $item->link . '',
                                            'date' =>   date('Y-m-d h:i:s A', strtotime($item->pubDate . '')),
                                            'description' => $item->description . '',
                                            'categories' => $categories,
                                            'author' => array( 'name' => $author),
                                            'content' => $content,

                                        );

                }
            }

            return $data;
        }

        /**
         * Parse feed xml into an array.
         *
         * @param object $feedXml SimpleXMLElementObject
         * @return array feed content
         * @access public
         */
        function parseAtom($feedXml) {
            $data = array();

            $data['title'] = $feedXml->title . '';
            foreach ($feedXml->link as $link) {
                    $data['link'] = $link['href'] . '';
                break;
            }

            $data['description'] = $feedXml->subtitle . '';
            $data['parser'] = __CLASS__;
            $data['type'] = 'atom';

            if (!empty($feedXml->entry)) {
                foreach ($feedXml->entry as $item) {
                    foreach ($item->link as $link) {
                        $itemLink = $link['href'] . '';
                        break;
                    }
                    $categories = array();
                    foreach ($item->category as $category) {
                        $categories[] = $category['term'] . '';
                    }

                    $data['items'][] = array(
                                            'title' =>  $item->title . '',
                                            'link' =>   $itemLink . '',
                                            'date' =>   date('Y-m-d h:i:s A', strtotime($item->published . '')),
                                            'description' => $item->summary . '',
                                            'content' => $item->content . '',
                                            'categories' => $categories,
                                            'author' => array('name' => $item->author->name . '', 'url' => $item->author->uri . ''),
                                            'extra' => array('contentType' => $item->content['type'] . '', 'descriptionType' => $item->summary['type'] . '')
                                        );
                }
            }

            return $data;
        }

    }
?>
