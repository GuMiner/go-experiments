using System;
using System.Collections.Generic;
using System.Drawing;
using System.IO;
using System.Linq;
using System.Windows.Forms;

namespace ReceiptEditor
{
    public partial class ReceiptEditor : Form, IDisposable
    {
        private const int MinSubImageSize = 10;

        private readonly Pen overlayPen = new Pen(Color.CornflowerBlue, 1);
        private readonly Pen highlightPen = new Pen(Color.LightGreen, 1);

        private Point lastMousePos = new Point(-1, -1);
        private Point lastClickedPos = new Point(-1, -1);
        private bool inSelectMode = false;

        private Queue<string> filesToProcess = new Queue<string>();
        private int filesProcessed = 0;

        private List<SubImage> subImages = new List<SubImage>();
        private List<IDisposable> forms = new List<IDisposable>();
        private string currentFileName = null;
        
        public ReceiptEditor()
        {
            InitializeComponent();

            ReceiptEditor.ImageCategories = ImageCategory.LoadImageCategories();
        }

        public static List<ImageCategory> ImageCategories { get; private set; }

        /// <summary>
        /// Render
        /// </summary>
        private void imageBox_Paint(object sender, PaintEventArgs e)
        {
            e.Graphics.DrawEllipse(overlayPen, new Rectangle(lastMousePos, new Size(10, 10)));

            if (inSelectMode)
            {
                e.Graphics.DrawRectangle(overlayPen, new Rectangle(lastClickedPos, new Size(lastMousePos.X - lastClickedPos.X, lastMousePos.Y - lastClickedPos.Y)));
            }

            foreach (SubImage subImage in subImages)
            {
                e.Graphics.DrawRectangle(highlightPen, new Rectangle(subImage.OriginalMin, new Size(subImage.OriginalMax.X - subImage.OriginalMin.X, subImage.OriginalMax.Y - subImage.OriginalMin.Y)));
            }
        }

        private void imageBox_MouseMove(object sender, MouseEventArgs e)
        {
            lastMousePos = e.Location;
            imageBox.Invalidate();
        }

        private void imageBox_MouseDown(object sender, MouseEventArgs e)
        {
            inSelectMode = true;
            lastClickedPos = e.Location;
        }

        private void imageBox_MouseUp(object sender, MouseEventArgs e)
        {
            if (Math.Abs(lastMousePos.X - lastClickedPos.X) < MinSubImageSize || Math.Abs(lastMousePos.Y - lastClickedPos.Y) < MinSubImageSize)
            {
                // Cancel the operation
                return;
            }
            else
            {
                SubImage subImage = new SubImage(GetImagePosition(lastClickedPos), GetImagePosition(lastMousePos), (Image)imageBox.Image.Clone(), lastClickedPos, lastMousePos);
                subImages.Add(subImage);

                Form2 form = new Form2(this.subImages.Count, subImage);
                form.Show();
                forms.Add(form);
            }

            inSelectMode = false;
        }

        /// <summary>
        /// Grabbed from https://www.codeproject.com/articles/20923/mouse-position-over-image-in-a-picturebox, with modifications (MIT)
        /// </summary>

        private Point GetImagePosition(Point coordinates)
        {
            // This is the one that gets a little tricky. Essentially, need to check the aspect ratio of the image to the aspect ratio of the control to determine how it is being rendered
            float imageAspect = (float)imageBox.Image.Width / imageBox.Image.Height;
            float controlAspect = (float)imageBox.Width / imageBox.Height;
            float newX = coordinates.X;
            float newY = coordinates.Y;
            if (imageAspect > controlAspect)
            {
                // This means that we are limited by width, meaning the image fills up the entire control from left to right
                float ratioWidth = (float)imageBox.Image.Width / imageBox.Width;
                newX *= ratioWidth;
                float scale = (float)imageBox.Width / imageBox.Image.Width;
                float displayHeight = scale * imageBox.Image.Height;
                float diffHeight = imageBox.Height - displayHeight;
                diffHeight /= 2;
                newY -= diffHeight;
                newY /= scale;
            }
            else
            {
                // This means that we are limited by height, meaning the image fills up the entire control from top to bottom
                float ratioHeight = (float)imageBox.Image.Height / imageBox.Height;
                newY *= ratioHeight;
                float scale = (float)imageBox.Height / imageBox.Image.Height;
                float displayWidth = scale * imageBox.Image.Width;
                float diffWidth = imageBox.Width - displayWidth;
                diffWidth /= 2;
                newX -= diffWidth;
                newX /= scale;
            }

            return new Point((int)newX, (int)newY);
        }

        private void ReceiptEditor_Resize(object sender, EventArgs e)
        {
            imageBox.Invalidate();
        }

        private void scanButton_Click(object sender, EventArgs e)
        {
            this.filesToProcess = this.FindFiles(this.receiptFolderBox.Text);
            this.IterateFiles(this.processedFolderBox.Text, (file) => ++this.filesProcessed);

            this.UpdatePercentDone();
            this.AdvanceImage(skipValidationSteps: true);
        }

        private void nextButton_Click(object sender, EventArgs e)
        {
            this.AdvanceImage();
        }

        private void AdvanceImage(bool skipValidationSteps = false)
        {
            // Validate we don't inadvertently lose data.
            if (!skipValidationSteps)
            {
                if (subImages.Count == 0)
                {
                    DialogResult result = MessageBox.Show("No subimages were selected. Do you want to continue?", "No Subimages selected", MessageBoxButtons.YesNo);
                    if (result != DialogResult.Yes)
                    {
                        return;
                    }
                }

                if (subImages.Any(image => !image.Saved))
                {
                    DialogResult result = MessageBox.Show("Not all subimages were saved. Do you want to continue?", "Not all subimages saved", MessageBoxButtons.YesNo);
                    if (result != DialogResult.Yes)
                    {
                        return;
                    }
                }
            }

            // Release all handles to the image
            this.subImages.Clear();
            foreach (IDisposable form in this.forms)
            {
                form.Dispose();
            }

            this.forms.Clear();

            Image mainImage = this.imageBox.Image;
            this.imageBox.Image = null;
            mainImage?.Dispose();

            // Move to the processed folder, if we actually have a file to update.
            if (this.currentFileName != null)
            {
                string destinationDirectory = Path.Combine(this.processedFolderBox.Text, Path.GetFileName(Path.GetDirectoryName(this.currentFileName)));
                string fileName = Path.GetFileName(this.currentFileName);
                string destinationFile = Path.Combine(destinationDirectory, fileName);
                while (File.Exists(destinationFile))
                {
                    fileName = $"0{fileName}";
                    destinationFile = Path.Combine(destinationDirectory, fileName);
                }

                if (!Directory.Exists(destinationDirectory))
                {
                    Directory.CreateDirectory(destinationDirectory);
                }

                File.Move(this.currentFileName, destinationFile);

                // Update UI.
                ++this.filesProcessed;
                this.UpdatePercentDone();
            }

            // Advance to the next image.
            this.currentFileName = this.filesToProcess.Dequeue();
            using (FileStream stream = new FileStream(this.currentFileName, FileMode.Open, FileAccess.Read))
            {
                this.imageBox.Image = Image.FromStream(stream);
            }

            imageFolderBox.Text = Path.GetDirectoryName(this.currentFileName);
            imageNameBox.Text = $"{Path.GetFileName(this.currentFileName)}: {this.imageBox.Image.Width}x{this.imageBox.Image.Height}";
        }

        private void UpdatePercentDone()
        {
            string percentDoneString = $"{(this.filesProcessed * 100) / (this.filesProcessed + this.filesToProcess.Count)}% done.";
            scanResultsTextBox.Text = $"{this.filesProcessed} proc., {this.filesToProcess.Count} ready, {percentDoneString}";
        }

        private Queue<string> FindFiles(string directory)
        {
            Queue<string> fileQueue = new Queue<string>();
            this.IterateFiles(directory, (file) => fileQueue.Enqueue(file));
            return fileQueue;
        }

        private void IterateFiles(string directory, Action<string> fileAction)
        {
            foreach (string file in Directory.GetFiles(directory))
            {
                fileAction(file);
            }

            foreach (string subDirectory in Directory.GetDirectories(directory))
            {
                IterateFiles(subDirectory, fileAction);
            }
        }
    }
}
