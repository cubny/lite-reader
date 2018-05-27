<?php 

namespace LiteReader\Unit\Feed\Domain\ValueObject;

use LiteReader\Feed\Domain\ValueObjects\PublishedAt;
use PHPUnit\Framework\TestCase;

final class PublishedAtTest extends TestCase
{
    public function testEmptyConstructorCreatesNow()
    {
        $publishedAt = new PublishedAt();
        $this->assertInstanceOf(\DateTime::class, $publishedAt->getValue());
    }
}
