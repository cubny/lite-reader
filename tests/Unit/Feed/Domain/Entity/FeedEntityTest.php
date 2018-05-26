<?php 

namespace LiteReader\Unit\Feed\Domain\Entity;

use LiteReader\Feed\Domain\Entity\FeedEntity;
use LiteReader\Feed\Domain\ValueObjects\{
    Id,
    Title,
    Description,
    Url,
    UpdatedAt
};
use PHPUnit\Framework\TestCase;

final class FeedEntityTest extends TestCase
{
    public function setUp()
    {
        $this->feed = new FeedEntity(
            new Id(1),
            new Title('title'),
            new Description('description description'),
            new Url('http://newfeed.com/rss')
        );
    }
    /**
     * @dataProvider attributeValueObjectItems
     */
    public function testAttributeIsInstanceOfValueObjects($type, $attribute)
    {
        $this->assertAttributeInstanceOf($type, $attribute, $this->feed);
    }

    public function attributeValueObjectItems()
    {
        return [
            "id" => [Id::class, 'id'],
            "title" => [Title::class, 'title'],
            "description" => [Description::class, 'description'],
            "Url" => [Url::class, 'url']
        ];
    }
}
