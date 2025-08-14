
import sys # Importing necessary libraries
import os # For file handling
import json # For JSON handling

try:
    from docx import Document #For parsing DOCX files
except ImportError:
    Document = None
try:
    import pdfplumber # For parsing PDF files
except ImportError:
    pdfplumber = None

def extract_text_from_pdf(file_path):# Function to extract text from PDF files
    if not pdfplumber:
        return None, "pdfplumber not installed" #return if not not installed
    try:
        text = ""
        with pdfplumber.open(file_path) as pdf:
            for page in pdf.pages:
                text += page.extract_text() or ""
        return text, None
    except Exception as e:
        return None, f"Error reading PDF: {str(e)}"

def extract_text_from_docx(file_path): # Function to extract text from DOCX files
    if not Document:
        return None, "python-docx not installed"
    try:
        doc = Document(file_path)
        return "\n".join([para.text for para in doc.paragraphs]), None
    except Exception as e:
        return None, f"Error reading DOCX: {str(e)}"

def main():
    if len(sys.argv) < 2:
        print(json.dumps({"error": "No file path provided"}))
        sys.exit(1)
    file_path = sys.argv[1]
    if not os.path.exists(file_path):
        print(json.dumps({"error": "File does not exist"}))
        sys.exit(1)
    ext = os.path.splitext(file_path)[1].lower()
    if ext == ".pdf":
        text, error = extract_text_from_pdf(file_path)
    elif ext == ".docx":
        text, error = extract_text_from_docx(file_path)
    else:
        print(json.dumps({"error": "Unsupported file type"}))
        sys.exit(1)
    if error:
        print(json.dumps({"error": error}))
        sys.exit(1)
    # Dummy parsing logic for now
    result = {
        "name": "John Doe",
        "email": "john.doe@example.com",
        "skills": ["Python", "Go", "Resume Parsing"],
        "raw_text": text[:500] if text else ""
    }
    print(json.dumps(result))

if __name__ == "__main__":
    main() 