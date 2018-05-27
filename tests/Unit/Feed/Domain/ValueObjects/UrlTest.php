<?php 

namespace LiteReader\Unit\Feed\Domain\ValueObject;

use InvalidArgumentException;
use LiteReader\Feed\Domain\ValueObjects\Url;
use PHPUnit\Framework\TestCase;

final class UrlTest extends TestCase
{
    public function testThrowsInvalidArgumentExceptionOnInvalidUrl()
    {
        $this->expectException(InvalidArgumentException::class);
        $url = new Url("com");
    }

    public function testConstructWithValidUrlCreatesUrl()
    {
        $url = new Url("http://localhost/rss");
        $this->assertInternalType("string", $url->getValue());
    }
}
