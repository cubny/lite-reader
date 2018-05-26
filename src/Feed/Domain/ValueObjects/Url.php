<?php

namespace LiteReader\Feed\Domain\ValueObjects;

final class Url
{
    /**
     *
     * @var string 
     */
    private $url;

    /**
     * Constructoror
     *
     * @param string $url
     * @access public
     * @return void
     */
    public function __constructor(string $url)
    {
        $this->url = $url;
    }

    /**
     * getValue
     *
     * @access public
     * @return string
     */
    public function getValue(): string
    {
        return $this->url;
    }
}
