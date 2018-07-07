using System;
using System.Drawing;
using System.Drawing.Imaging;
using System.Windows.Forms;

namespace ReceiptEditor
{
    public partial class Form2 : Form, IDisposable
    {
        private readonly int subImageId;
        private readonly SubImage subImage;
        private ImageAttributes imageAttributes;

        private bool inTranslateMode = false;
        private Point lastMousePos = new Point(-1, -1);

        private ImageEditForm imageEditForm;

        public Form2()
        {
            InitializeComponent();
            this.imageAttributes = new ImageAttributes();
            this.MouseWheel += new MouseEventHandler(MouseWheelMove);
        }

        private void MouseWheelMove(object sender, MouseEventArgs e)
        {
            float aspectRatio = 1.0f / ((float)subImage.Image.Width / (float)subImage.Image.Height);
            int scrollDelta = e.Delta;
            int zoomAddFactor = 10;
            if (scrollDelta > 0)
            {
                subImage.MinPos = new Point((int)((float)subImage.MinPos.X - zoomAddFactor * aspectRatio), (int)((float)subImage.MinPos.Y - zoomAddFactor));
                subImage.MaxPos = new Point((int)((float)subImage.MaxPos.X + zoomAddFactor * aspectRatio), (int)((float)subImage.MaxPos.Y + zoomAddFactor));
            }
            else
            {
                subImage.MinPos = new Point((int)((float)subImage.MinPos.X + zoomAddFactor* aspectRatio), (int)((float)subImage.MinPos.Y + zoomAddFactor));
                subImage.MaxPos = new Point((int)((float)subImage.MaxPos.X - zoomAddFactor* aspectRatio), (int)((float)subImage.MaxPos.Y - zoomAddFactor));
            }

            imageBox.Invalidate();
        }

        public Form2(int subImageId, SubImage subImage)
            : this()
        {
            this.Text = $"Sub Image - {subImageId}";
            this.subImageId = subImageId;
            this.subImage = subImage;

            this.imageEditForm = new ImageEditForm(
                this.subImageId,
                (attr) =>
                {
                    this.imageAttributes = attr;
                    imageBox.Invalidate();
                },
                () => {
                    this.subImage.Saved = true;
                    int width = subImage.MaxPos.X - subImage.MinPos.X;
                    int height = subImage.MaxPos.Y - subImage.MinPos.Y;
                    Bitmap bitmap = new Bitmap(width, height);
                    using (Graphics g = Graphics.FromImage(bitmap))
                    {
                        DrawPartialImage(g, width, height);
                    }

                    return bitmap;
                }, () => this.Hide());

            this.Show();
            this.imageEditForm.Show();
        }

        public new void Dispose()
        {
            this.imageBox.Image = null;
            this.imageEditForm.Dispose();
            base.Dispose();
        }

        private void pictureBox1_MouseMove(object sender, MouseEventArgs e)
        {
            if (inTranslateMode)
            {
                int deltaX = e.X - lastMousePos.X;
                int deltaY = e.Y - lastMousePos.Y;

                float xScaleFactor = (float)(subImage.MaxPos.X - subImage.MinPos.X) / (float)imageBox.Width;
                float yScaleFactor = (float)(subImage.MaxPos.Y - subImage.MinPos.Y) / (float)imageBox.Height;
                subImage.MinPos = new Point((int)((float)subImage.MinPos.X - (float)deltaX * xScaleFactor), (int)((float)subImage.MinPos.Y - (float)deltaY * yScaleFactor));
                subImage.MaxPos = new Point((int)((float)subImage.MaxPos.X - (float)deltaX * xScaleFactor), (int)((float)subImage.MaxPos.Y - (float)deltaY * yScaleFactor));

                lastMousePos = e.Location;
                imageBox.Invalidate();
            }
        }

        private void pictureBox1_MouseDown(object sender, MouseEventArgs e)
        {
            inTranslateMode = true;
            lastMousePos = e.Location;
        }

        private void pictureBox1_MouseUp(object sender, MouseEventArgs e)
        {
            inTranslateMode = false;
        }

        private void imageBox_Paint(object sender, PaintEventArgs e)
        {
            DrawPartialImage(e.Graphics, imageBox.Width, imageBox.Height);
        }

        private void DrawPartialImage(Graphics g, int destinationWidth, int destinationHeight)
        {
            g.Clear(Color.LightSeaGreen);
            g.DrawImage(subImage.Image,
                new Rectangle(0, 0, destinationWidth, destinationHeight),
                subImage.MinPos.X, subImage.MinPos.Y,
                subImage.MaxPos.X - subImage.MinPos.X,
                subImage.MaxPos.Y - subImage.MinPos.Y,
                GraphicsUnit.Pixel,
                imageAttributes);
        }

        private void Form2_Resize(object sender, EventArgs e)
        {
            imageBox.Invalidate();
        }
    }
}
