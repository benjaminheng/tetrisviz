# Examples
<table>
<tr><th>tetrisviz</th><th>pikchr</th><th>svg</th></tr>
<tr>
<td><pre>
+board 10x20

-r---t-bbb
rrg--t-yyb
rpgg-t-yyo
pppg-t-ooo

</pre></td>
<td><pre>
boxwid = 0.2
boxht = boxwid
$currLine = 1
define next {
  box invis at (-boxwid, -boxwid*$currLine)
  $currLine = $currLine + 1
}
define $e { box fill 0xc1c1c1 }
define $r { box fill 0xf13636 }
define $t { box fill 0x67edf5 }
define $b { box fill skyblue }
define $g { box fill 0x39e572 }
define $y { box fill 0xfff223 }
define $p { box fill 0xc936f1 }
define $o { box fill 0xfbbb11 }

$e;$e;$e;$e;$e;$e;$e;$e;$e;$e;next;
$e;$e;$e;$e;$e;$e;$e;$e;$e;$e;next;
$e;$e;$e;$e;$e;$e;$e;$e;$e;$e;next;
$e;$e;$e;$e;$e;$e;$e;$e;$e;$e;next;
$e;$e;$e;$e;$e;$e;$e;$e;$e;$e;next;
$e;$e;$e;$e;$e;$e;$e;$e;$e;$e;next;
$e;$e;$e;$e;$e;$e;$e;$e;$e;$e;next;
$e;$e;$e;$e;$e;$e;$e;$e;$e;$e;next;
$e;$e;$e;$e;$e;$e;$e;$e;$e;$e;next;
$e;$e;$e;$e;$e;$e;$e;$e;$e;$e;next;
$e;$e;$e;$e;$e;$e;$e;$e;$e;$e;next;
$e;$e;$e;$e;$e;$e;$e;$e;$e;$e;next;
$e;$e;$e;$e;$e;$e;$e;$e;$e;$e;next;
$e;$e;$e;$e;$e;$e;$e;$e;$e;$e;next;
$e;$e;$e;$e;$e;$e;$e;$e;$e;$e;next;
$e;$e;$e;$e;$e;$e;$e;$e;$e;$e;next;
$e;$r;$e;$e;$e;$t;$e;$b;$b;$b;next;
$r;$r;$g;$e;$e;$t;$e;$y;$y;$b;next;
$r;$p;$g;$g;$e;$t;$e;$y;$y;$o;next;
$p;$p;$p;$g;$e;$t;$e;$o;$o;$o;next;$e;$t;$e;$y;$y;$o;next;
$p;$p;$p;$g;$e;$t;$e;$o;$o;$o;next;
</pre></td>
<td><img src="example1.svg" style="width: 200px"/></td>
</tr>
</table>