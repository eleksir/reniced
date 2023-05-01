## no critic
# Почему-то perlcritic пытается критиковать и cpanfile тоже, хотя он к перловым программам и либам не относится.
# Проставим версии зависимостей в то, что по идее должно работать с текущим перлом. Типа "как у взрослых".
requires Encode,       '>=3.19';
requires Getopt::Long, '>=2.54';
requires JSON::XS,     '>=4.03';
# Built-in perl distribution
requires POSIX,        '0';
requires Proc::Find,   '>=0.051';
requires Time::HiRes,  '>=1.9764';
