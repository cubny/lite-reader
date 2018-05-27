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
    public function __construct(string $url)
    {
        $url = filter_var($url, FILTER_SANITIZE_URL);
        if (\filter_var($url, FILTER_VALIDATE_URL, FILTER_FLAG_PATH_REQUIRED) === false) {
            throw new \InvalidArgumentException();
        }
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
