using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.Data;
using System.Drawing;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using System.Windows.Forms;

namespace ReceiptEditor
{
    public partial class ReceiptEditor : Form, IDisposable
    {
        private Pen overlayPen = new Pen(Color.CornflowerBlue, 1);
        private Pen highlightPen = new Pen(Color.LightGreen, 1);

        private Point lastMousePos = new Point(-1, -1);
        private Point lastClickedPos = new Point(-1, -1);
        private bool inSelectMode = false;
        private int minSubImageSize = 10;

        private List<SubImage> subImages = new List<SubImage>();

        public ReceiptEditor()
        {
            InitializeComponent();
            this.imageBox.Image = Image.FromFile(@"C:\Users\Gustave\Desktop\Data Archive\receipts\2015-01-21\001.jpg");
        }

        // Unused?
        private void Form1_Load(object sender, EventArgs e)
        {

        }

        private void nextButton_Click(object sender, EventArgs e)
        {

        }

        // Unused
        private void label2_Click(object sender, EventArgs e)
        {
            
        }

        private void textBox1_TextChanged(object sender, EventArgs e)
        {

        }

        private void categoryBox_SelectedIndexChanged(object sender, EventArgs e)
        {

        }

        private void itemDate_ValueChanged(object sender, EventArgs e)
        {

        }

        private void imageBox_Click(object sender, EventArgs e)
        {

        }

        /// <summary>
        /// Render
        /// </summary>
        private void imageBox_Paint(object sender, PaintEventArgs e)
        {
            // e.Graphics.DrawEllipse(overlayPen, 0, 0, this.imageBox.Width, this.imageBox.Height);
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
            if (lastMousePos.X < lastClickedPos.X + minSubImageSize || lastMousePos.Y < lastClickedPos.Y + minSubImageSize)
            {
                // Cancel the operation
                return;
            }
            else
            {
                SubImage subImage = new SubImage(GetImagePosition(lastClickedPos), GetImagePosition(lastMousePos), (Image)imageBox.Image.Clone(), lastClickedPos, lastMousePos);
                subImages.Add(subImage);

                Form2 form = new Form2(subImage);
                form.Show();
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

        private void button1_Click(object sender, EventArgs e)
        {

        }

        private void addCategoryField_TextChanged(object sender, EventArgs e)
        {

        }

        private void ReceiptEditor_Resize(object sender, EventArgs e)
        {
            imageBox.Invalidate();
        }
    }
}
