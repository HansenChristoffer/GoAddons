# This file is part of GoAddons, which is licensed under the GNU General Public License v3.0.
# You should have received a copy of the GNU General Public License along with this program.
# If not, see <https://www.gnu.org/licenses/>.

CREATE DATABASE IF NOT EXISTS `defcon` CHARSET=latin1 COLLATE=latin1_swedish_ci;

CREATE TABLE IF NOT EXISTS `defcon`.`system_config_default` (
  `name` VARCHAR(255),
  `path` VARCHAR(255),
  PRIMARY KEY (`name`)
) CHARSET=latin1 COLLATE=latin1_swedish_ci COMMENT 'Table with the systems default configurations. E.g. path to where to download files';

REPLACE INTO `defcon`.`system_config_default` (`name`, `path`)
  VALUES ('browser.download.dir', '/home/miso/Downloads/kaasufouji-addons');

REPLACE INTO `defcon`.`system_config_default` (`name`, `path`)
  VALUES ('extract.addon.path', '/home/miso/ssd/Games/battlenet/drive_c/Program Files (x86)/World of Warcraft/_retail_/Interface/AddOns');

CREATE DATABASE IF NOT EXISTS `kaasufouji` CHARSET=latin1 COLLATE=latin1_swedish_ci;

CREATE TABLE IF NOT EXISTS `kaasufouji`.`addons` (
  `id` INT AUTO_INCREMENT PRIMARY KEY,
  `name` VARCHAR(128) NOT NULL UNIQUE,
  `filename` VARCHAR(255) DEFAULT NULL,
  `url` VARCHAR(255) NULL DEFAULT NULL,
  `download_url` VARCHAR(255) NULL DEFAULT NULL,
  `update_available` BOOLEAN NULL DEFAULT FALSE,
  `last_downloaded` TIMESTAMP NULL DEFAULT NULL,
  `last_modified_at` TIMESTAMP DEFAULT (CURRENT_TIMESTAMP),
  `added_at` TIMESTAMP DEFAULT (CURRENT_TIMESTAMP)
) CHARSET=latin1 COLLATE=latin1_swedish_ci COMMENT 'Table consists of Addons and where to find them';

CREATE TABLE IF NOT EXISTS `kaasufouji`.`run_log` (
  `id` INT AUTO_INCREMENT PRIMARY KEY,
  `run_id` VARCHAR(255) NOT NULL,
  `service` VARCHAR(128) NOT NULL,
  `added_at` TIMESTAMP DEFAULT (CURRENT_TIMESTAMP)
) CHARSET=latin1 COLLATE=latin1_swedish_ci COMMENT 'Log table for all service running';

CREATE TABLE IF NOT EXISTS `kaasufouji`.`download_log` (
  `id` INT AUTO_INCREMENT PRIMARY KEY,
  `run_id` VARCHAR(255) NOT NULL,
  `url` VARCHAR(255) NOT NULL,
  `added_at` TIMESTAMP DEFAULT (CURRENT_TIMESTAMP)
) CHARSET=latin1 COLLATE=latin1_swedish_ci COMMENT 'Log table for all downloads';

CREATE TABLE IF NOT EXISTS `kaasufouji`.`extract_log` (
  `id` INT AUTO_INCREMENT PRIMARY KEY,
  `run_id` VARCHAR(255) NOT NULL,
  `file` TEXT NOT NULL,
  `added_at` TIMESTAMP DEFAULT (CURRENT_TIMESTAMP)
) CHARSET=latin1 COLLATE=latin1_swedish_ci COMMENT 'Log table for all addon extractions';

CREATE USER IF NOT EXISTS 'tanushi'@'%' IDENTIFIED BY 'shitanu';
GRANT SELECT, INSERT, UPDATE ON kaasufouji.* TO 'tanushi'@'%';
GRANT SELECT ON defcon.* TO 'tanushi'@'%';
FLUSH PRIVILEGES;
