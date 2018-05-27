<?php 

namespace LiteReader\Unit\Feed\Domain\ValueObject;

use LiteReader\Feed\Domain\ValueObjects\UpdatedAt;
use PHPUnit\Framework\TestCase;

final class UpdatedAtTest extends TestCase
{
    public function testEmptyConstructorCreatesNow()
    {
        $publishedAt = new UpdatedAt();
        $this->assertInstanceOf(\DateTime::class, $publishedAt->getValue());
    }
}
