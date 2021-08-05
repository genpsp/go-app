-- +migrate Up
CREATE TABLE `item` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `user_id` BIGINT UNSIGNED NOT NULL,
    `created_at` DATETIME NOT NULL,
    `updated_at` DATETIME NULL,
    `deleted_at` DATETIME NULL,
    PRIMARY KEY (`id`),
    INDEX `fk_item_user1_idx` (`user_id` ASC),
    CONSTRAINT `fk_item_user1`
    FOREIGN KEY (`user_id`)
     REFERENCES `user` (`id`)
     ON DELETE NO ACTION
     ON UPDATE NO ACTION)
ENGINE = InnoDB;


-- +migrate Down
DROP TABLE `item`;
