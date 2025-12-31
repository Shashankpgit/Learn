Here are **quick, clean notes** on **creating files and file-related commands in Linux/Unix**, from **basic → advanced**.
(Perfect for revision or interviews.)

---

## 1️⃣ Creating Files

### Basic Ways

```bash
touch file.txt            # Create empty file
cat > file.txt             # Create + write (Ctrl+D to save)
echo "Hello" > file.txt    # Create + overwrite
echo "Hello" >> file.txt   # Append
```

### Using Editors

```bash
nano file.txt
vi file.txt
vim file.txt
```

---

## 2️⃣ Viewing Files

```bash
cat file.txt               # Show full file
less file.txt              # Scrollable view
more file.txt
head file.txt              # First 10 lines
tail file.txt              # Last 10 lines
tail -f logfile.log        # Live monitoring
```

---

## 3️⃣ File Operations

```bash
cp a.txt b.txt              # Copy
mv a.txt b.txt              # Rename / move
rm file.txt                 # Delete
rm -r folder/               # Delete directory
```

⚠️ **Dangerous**

```bash
rm -rf /folder              # Force delete (no recovery)
```

---

## 4️⃣ File Permissions

```bash
ls -l file.txt              # View permissions
chmod 777 file.txt          # Full permission
chmod u+x file.sh           # Add execute
chown user:group file.txt   # Change owner
```

Permission format:

```
-rwxr-xr--
```

---

## 5️⃣ File Search & Info

```bash
find / -name file.txt       # Search file
locate file.txt             # Fast search
file file.txt               # File type
stat file.txt               # Detailed info
```

---

## 6️⃣ Text Processing (Intermediate)

```bash
grep "word" file.txt        # Search text
wc file.txt                 # Line, word, char count
sort file.txt               # Sort lines
uniq file.txt               # Remove duplicates
cut -d: -f1 file.txt        # Column extraction
```

---

## 7️⃣ Redirection & Pipes

```bash
command > file.txt          # Output overwrite
command >> file.txt         # Append
command < file.txt          # Input
command1 | command2         # Pipe output
```

Example:

```bash
ls | grep ".txt"
```

---

## 8️⃣ Compression & Archiving

```bash
tar -cvf files.tar file/
tar -xvf files.tar
gzip file.txt
gunzip file.txt.gz
zip file.zip file.txt
unzip file.zip
```

---

## 9️⃣ Advanced File Commands

```bash
ln file.txt link.txt        # Hard link
ln -s file.txt slink.txt    # Symbolic link
xargs                       # Build arguments
awk '{print $1}' file.txt  # Text processing
sed 's/old/new/g' file.txt # Replace text
```

---

## 10️⃣ Permissions + Execution

```bash
chmod +x script.sh
./script.sh
```

---

### ✅ Pro Tip (Most Used Combo)

```bash
ls | grep name | wc -l
```

---
Below is a **clear explanation of Linux file commands with real, practical examples**, starting from **basic → advanced**, exactly how you’d use them on a real system.

---

## 1️⃣ Creating Files (Real Use)

### `touch`

👉 Create an empty file (often for logs or placeholders)

```bash
touch report.txt
```

✔ Creates `report.txt` if it doesn’t exist
✔ Updates timestamp if it exists

---

### `cat`

👉 Create and write content

```bash
cat > notes.txt
This is Linux
File handling basics
Ctrl+D
```

✔ Saves text after `Ctrl+D`

---

### `echo`

👉 Quick content creation

```bash
echo "Server started" > log.txt
```

✔ Overwrites file

```bash
echo "New entry" >> log.txt
```

✔ Appends content

---

## 2️⃣ Viewing Files

### `cat`

👉 View small files

```bash
cat notes.txt
```

---

### `less`

👉 Best for large files

```bash
less /var/log/syslog
```

Controls:

* `Enter` → line
* `Space` → page
* `q` → quit

---

### `head` / `tail`

👉 Check beginning or end of files

```bash
head -5 users.txt
tail -10 errors.log
```

👉 Live log monitoring

```bash
tail -f app.log
```

---

## 3️⃣ Copy, Move, Delete

### `cp`

```bash
cp file1.txt backup.txt
cp -r project/ project_backup/
```

---

### `mv`

```bash
mv old.txt new.txt
mv report.txt /home/user/docs/
```

---

### `rm`

```bash
rm temp.txt
rm -r old_folder/
```

⚠️ **No recycle bin**

```bash
rm -rf folder/
```

---

## 4️⃣ File Permissions (Very Important)

### View permissions

```bash
ls -l script.sh
```

Output:

```
-rw-r--r--
```

### Change permissions

```bash
chmod 644 file.txt
chmod u+x script.sh
```

### Change ownership

```bash
chown john:dev file.txt
```

---

## 5️⃣ Searching Files

### `find`

```bash
find /home -name resume.pdf
```

### `locate`

```bash
locate config.yml
```

✔ Faster but needs updated database

---

### File info

```bash
file image.png
stat file.txt
```

---

## 6️⃣ Text Search & Analysis

### `grep`

```bash
grep "error" app.log
grep -i "fail" report.txt
grep -r "TODO" project/
```

---

### `wc`

```bash
wc notes.txt
```

Output:

```
Lines Words Characters
```

---

### `sort` & `uniq`

```bash
sort names.txt
sort names.txt | uniq
```

---

### `cut`

```bash
cut -d: -f1 /etc/passwd
```

✔ Extract usernames

---

## 7️⃣ Redirection & Pipes (Power of Linux)

### Redirect output

```bash
ls > files.txt
```

### Pipe commands

```bash
ls | grep ".txt"
```

### Combine commands

```bash
cat access.log | grep 404 | wc -l
```

✔ Count errors

---

## 8️⃣ Compression & Archives

### `tar`

```bash
tar -cvf backup.tar project/
tar -xvf backup.tar
```

### `gzip`

```bash
gzip largefile.txt
gunzip largefile.txt.gz
```

### `zip`

```bash
zip files.zip a.txt b.txt
unzip files.zip
```

---

## 9️⃣ Advanced Commands (Real Power)

### Links

```bash
ln file.txt hardlink.txt
ln -s file.txt softlink.txt
```

---

### `sed`

👉 Replace text

```bash
sed 's/dev/prod/g' config.txt
```

---

### `awk`

👉 Column-based processing

```bash
awk '{print $1}' users.txt
```

---

### `xargs`

```bash
cat files.txt | xargs rm
```

---

## 🔟 Running Scripts

```bash
chmod +x deploy.sh
./deploy.sh
```

---

## ✅ Real-World Combo Examples

✔ Find large files:

```bash
find / -size +100M
```

✔ Kill process by name:

```bash
ps aux | grep chrome | awk '{print $2}' | xargs kill
```

✔ Count `.log` files:

```bash
ls *.log | wc -l
```

---


