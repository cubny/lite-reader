<?php

/*
   Script: Helper/Date.php
      Display the number of minutes, hours, days, weeks or month ago.

   License:
      http://www.opensource.org/licenses/mit-license.html MIT License

   Copyright:
      Philippe Archambault <philippe.archambault@gmail.com>
*/

/*
   Function: pretty_date
   Display the number of minuts, hours, days, weeks or month ago.
*/

function pretty_date($from_time, $to_time = null)
{
    $to_time = $to_time ? $to_time: $_SERVER['REQUEST_TIME'];

    $distance_in_minutes = floor(abs($to_time - $from_time) / 60);

    if ($distance_in_minutes <= 1)
        return 'less then a minute';
    else if ($distance_in_minutes < 60)
        return $distance_in_minutes . ' minutes ago';
    else if ($distance_in_minutes < 90)
        return '1 hour ago';
    else if ($distance_in_minutes < 1440)
        return round($distance_in_minutes / 60) . ' hours ago', );
    else if ($distance_in_minutes < 2880)
        return 'Yesterday';
    else if ($distance_in_minutes < 10080)
        return round($distance_in_minutes / 1440) . ' days ago';
    else if ($distance_in_minutes < 43200)
        return round($distance_in_minutes / 10080) . ' weeks ago';
    else if ($distance_in_minutes < 86400)
        return '1 month ago';
    else if ($distance_in_minutes < 525960)
        return round($distance_in_minutes / 43200) . ' months ago';
    else if ($distance_in_minutes < 1051920)
        return '1 year ago';
    else
        return 'more then ' . round($distance_in_minutes / 525960) .' years ago';
}
